package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	keyfunc "github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const contextClaimsKey = "claims"
const contextUserIDKey = "user_id"
const contextTokenKey = "token"
const contextSessionIDKey = "session_id"

// ---------- JWKS cache (ต่อ "ชุด" ของ URL) ----------
var (
	jwksMu   sync.RWMutex
	jwksPool = map[string]keyfunc.Keyfunc{}
)

func canonicalJWKSKey(urls []string) string {
	// normalize, sort แล้ว join เพื่อให้ key เดิมทุกครั้ง
	norm := make([]string, 0, len(urls))
	for _, u := range urls {
		s := strings.TrimSpace(u)
		if s == "" {
			continue
		}
		// ไม่ตัด slash ท้ายเพื่อหลีกเลี่ยงชนกันกับ endpoint อื่น
		norm = append(norm, s)
	}
	sort.Strings(norm)
	return strings.Join(norm, " | ")
}

func getKeyfunc(jwksURLs []string) (keyfunc.Keyfunc, error) {
	if len(jwksURLs) == 0 {
		return nil, fmt.Errorf("jwksURLs is required")
	}
	key := canonicalJWKSKey(jwksURLs)

	// cache hit
	jwksMu.RLock()
	k := jwksPool[key]
	jwksMu.RUnlock()
	if k != nil {
		return k, nil
	}

	// สร้างครั้งเดียว พร้อม auto-refresh (v3)
	kk, err := keyfunc.NewDefaultCtx(context.Background(), jwksURLs)
	if err != nil {
		return nil, fmt.Errorf("jwks init failed: %w", err)
	}

	jwksMu.Lock()
	jwksPool[key] = kk
	jwksMu.Unlock()
	return kk, nil
}

func bearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("Authorization header must start with Bearer")
	}
	return parts[1], nil
}

func deriveIssuerFromJWKSs(jwksURLs []string) string {
	// บังคับ issuer ได้เฉพาะกรณีมี URL เดียวและเป็นรูป Keycloak มาตรฐาน
	if len(jwksURLs) != 1 {
		return ""
	}
	u := strings.TrimRight(jwksURLs[0], "/")
	const suf = "/protocol/openid-connect/certs"
	if strings.HasSuffix(u, suf) {
		return strings.TrimSuffix(u, suf)
	}
	return ""
}

// Helper function สำหรับสร้าง JWT parser options
func (m *GoMiddleware) buildJWTParserOptions(includeAudience bool) []jwt.ParserOption {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods(m.allowedAlgs),
		jwt.WithLeeway(30 * time.Second),
	}

	issuer := deriveIssuerFromJWKSs(m.jwksURL)
	if issuer != "" {
		opts = append(opts, jwt.WithIssuer(issuer))
	}

	if includeAudience && len(m.audience) > 0 {
		// กรอง audience ที่ไม่ว่างเปล่า
		validAudiences := make([]string, 0, len(m.audience))
		for _, aud := range m.audience {
			if trimmed := strings.TrimSpace(aud); trimmed != "" {
				validAudiences = append(validAudiences, trimmed)
			}
		}
		if len(validAudiences) > 0 {
			opts = append(opts, jwt.WithAudience(validAudiences...))
		}
	}

	return opts
}

// Helper function สำหรับ parse JWT token ด้วย JWKS
func (m *GoMiddleware) parseJWTWithJWKS(c echo.Context, tokenString string, includeAudience bool) (*jwt.Token, error) {
	kf, err := getKeyfunc(m.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("JWKS init failure: %w", err)
	}

	opts := m.buildJWTParserOptions(includeAudience)
	parser := jwt.NewParser(opts...)
	keyfuncCtx := kf.KeyfuncCtx(c.Request().Context())

	token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// บังคับต้องมี kid กัน key confusion
		if _, ok := t.Header["kid"]; !ok {
			return nil, errors.New("missing kid")
		}
		return keyfuncCtx(t)
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}

// Helper function สำหรับ parse JWT token ด้วย HMAC secret
func (m *GoMiddleware) parseJWTWithSecret(tokenString string, includeAudience bool) (*jwt.Token, error) {
	opts := []jwt.ParserOption{
		jwt.WithValidMethods(m.allowedAlgs),
		jwt.WithLeeway(30 * time.Second),
	}

	if iss := strings.TrimSpace(m.issuer); iss != "" {
		opts = append(opts, jwt.WithIssuer(iss))
	}

	if includeAudience && len(m.audience) > 0 {
		// กรอง audience ที่ไม่ว่างเปล่า
		validAudiences := make([]string, 0, len(m.audience))
		for _, aud := range m.audience {
			if trimmed := strings.TrimSpace(aud); trimmed != "" {
				validAudiences = append(validAudiences, trimmed)
			}
		}
		if len(validAudiences) > 0 {
			opts = append(opts, jwt.WithAudience(validAudiences...))
		}
	}

	parser := jwt.NewParser(opts...)

	token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// ยืนยันว่าเป็น HMAC จริง ๆ (hardening)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}

// Helper function สำหรับ parse JWT token แบบไม่ validate claims (สำหรับ SessionMiddleware)
func (m *GoMiddleware) parseJWTWithoutClaimsValidation(c echo.Context, tokenString string) (*jwt.Token, error) {
	kf, err := getKeyfunc(m.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("JWKS init failure: %w", err)
	}

	issuer := deriveIssuerFromJWKSs(m.jwksURL)
	opts := []jwt.ParserOption{
		jwt.WithValidMethods(m.allowedAlgs),
		jwt.WithoutClaimsValidation(), // ไม่ตรวจสอบ exp, iat, nbf
	}
	if issuer != "" {
		opts = append(opts, jwt.WithIssuer(issuer))
	}

	parser := jwt.NewParser(opts...)
	keyfuncCtx := kf.KeyfuncCtx(c.Request().Context())

	token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Header["kid"]; !ok {
			return nil, errors.New("missing kid")
		}
		return keyfuncCtx(t)
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}

// Helper function สำหรับ set claims ใน context
func setTokenClaims(c echo.Context, token *jwt.Token, tokenString string) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set(contextClaimsKey, claims)
		if userID, err := claims.GetSubject(); err == nil {
			c.Set(contextUserIDKey, userID)
		}
		c.Set(contextTokenKey, tokenString)
	}
}

// -------------------------------------------------------------

// IsAuthorizationWithKeycloak: ตรวจ JWT ด้วย JWKS ของ Keycloak
func (m *GoMiddleware) IsAuthorizationWithKeycloak(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := bearerToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		token, err := m.parseJWTWithJWKS(c, raw, true)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		setTokenClaims(c, token, raw)
		return next(c)
	}
}

// IsAuthorizationBasicJwt: ตรวจแบบ “secret” พื้นฐาน (HMAC) — ใช้เมื่อคุณควบคุมผู้ออกเอง
func (m *GoMiddleware) IsAuthorizationBasicJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		raw, err := bearerToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		token, err := m.parseJWTWithSecret(raw, true)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		setTokenClaims(c, token, raw)
		return next(c)
	}
}

// Role: ผ่านเมื่อผู้ใช้มีอย่างน้อย 1 role ตรงกับ requiredRoles
// รองรับรูปแบบ: "admin", "auth-service:admin", "realm:admin"
func (m *GoMiddleware) Role(requiredRoles []string) echo.MiddlewareFunc {
	req := make(map[string]struct{}, len(requiredRoles))
	for _, r := range requiredRoles {
		req[strings.ToLower(strings.TrimSpace(r))] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := c.Get(contextClaimsKey)
			claims, ok := val.(jwt.MapClaims)
			if !ok || claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No user claims found. Please authenticate first"})
			}
			userRoles := collectRolesKeycloak(claims)
			if len(userRoles) == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Insufficient permissions. Required roles: " + strings.Join(requiredRoles, ", ")})
			}

			// สร้าง role set จาก user roles
			roleSet := map[string]struct{}{}
			for _, r := range userRoles {
				normalizedRole := strings.ToLower(strings.TrimSpace(r))
				roleSet[normalizedRole] = struct{}{}

				// เพิ่มรูปแบบ shorthand สำหรับ resource roles ("admin" จาก "auth-service:admin")
				if strings.Contains(normalizedRole, ":") {
					parts := strings.SplitN(normalizedRole, ":", 2)
					if len(parts) == 2 {
						roleSet[parts[1]] = struct{}{} // เพิ่มเฉพาะชื่อ role โดยไม่มี prefix
					}
				}
			}

			// เช็ค any-of: มี role ใดตรงบ้าง
			for reqRole := range req {
				if _, ok := roleSet[reqRole]; ok {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": fmt.Sprintf("Insufficient permissions. Required roles: %v, User roles: %v", requiredRoles, userRoles),
			})
		}
	}
}

// ดึง roles
func collectRolesKeycloak(claims jwt.MapClaims) []string {
	var out []string

	// realm_access.roles
	if raRaw, ok := claims["realm_access"].(map[string]any); ok {
		if roles, ok := raRaw["roles"].([]any); ok {
			for _, r := range roles {
				if s, ok := r.(string); ok {
					out = append(out, "realm:"+s) // เพิ่ม prefix เพื่อแยกจาก resource roles
				}
			}
		}
	}

	// resource_access.{client_id}.roles - แยกตาม client
	if rsrc, ok := claims["resource_access"].(map[string]any); ok {
		for clientID, v := range rsrc {
			if m, ok := v.(map[string]any); ok {
				if roles, ok := m["roles"].([]any); ok {
					for _, r := range roles {
						if s, ok := r.(string); ok {
							out = append(out, clientID+":"+s) // รูปแบบ client_id:role_name
						}
					}
				}
			}
		}
	}

	// roles[] (fallback สำหรับ token รูปแบบเก่า)
	if arr, ok := claims["roles"].([]any); ok {
		for _, r := range arr {
			if s, ok := r.(string); ok {
				out = append(out, s)
			}
		}
	}

	// unique
	seen := map[string]struct{}{}
	uniq := make([]string, 0, len(out))
	for _, r := range out {
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		uniq = append(uniq, r)
	}
	return uniq
}

type Permission struct {
	ResourceID   string   `json:"rsid"`
	ResourceName string   `json:"rsname"`
	Scopes       []string `json:"scopes"`
}

// collectPermissions ดึง permissions จาก authorization.permissions
func collectPermissions(claims jwt.MapClaims) []Permission {
	var permissions []Permission

	if authz, ok := claims["authorization"].(map[string]any); ok {
		if perms, ok := authz["permissions"].([]any); ok {
			for _, perm := range perms {
				if permMap, ok := perm.(map[string]any); ok {
					permission := Permission{}

					if rsid, ok := permMap["rsid"].(string); ok {
						permission.ResourceID = rsid
					}
					if rsname, ok := permMap["rsname"].(string); ok {
						permission.ResourceName = rsname
					}
					if scopes, ok := permMap["scopes"].([]any); ok {
						for _, scope := range scopes {
							if scopeStr, ok := scope.(string); ok {
								permission.Scopes = append(permission.Scopes, scopeStr)
							}
						}
					}
					permissions = append(permissions, permission)
				}
			}
		}
	}

	return permissions
}

// Helper functions สำหรับดึงข้อมูลจาก context
func GetUserRoles(c echo.Context) []string {
	val := c.Get(contextClaimsKey)
	claims, ok := val.(jwt.MapClaims)
	if !ok || claims == nil {
		return []string{}
	}
	return collectRolesKeycloak(claims)
}

func GetUserPermissions(c echo.Context) []Permission {
	val := c.Get(contextClaimsKey)
	claims, ok := val.(jwt.MapClaims)
	if !ok || claims == nil {
		return []Permission{}
	}
	return collectPermissions(claims)
}

func GetUserID(c echo.Context) string {
	val := c.Get(contextUserIDKey)
	if userID, ok := val.(string); ok {
		return userID
	}
	return ""
}

func GetSessionID(c echo.Context) string {
	val := c.Get(contextSessionIDKey)
	if sessionID, ok := val.(string); ok {
		return sessionID
	}
	return ""
}

// HasRole ตรวจสอบว่าผู้ใช้มี role ที่ระบุหรือไม่
func HasRole(c echo.Context, role string) bool {
	userRoles := GetUserRoles(c)
	normalizedRole := strings.ToLower(strings.TrimSpace(role))

	for _, r := range userRoles {
		normalizedUserRole := strings.ToLower(strings.TrimSpace(r))
		if normalizedUserRole == normalizedRole {
			return true
		}
		// เช็ครูปแบบ shorthand
		if strings.Contains(normalizedUserRole, ":") {
			parts := strings.SplitN(normalizedUserRole, ":", 2)
			if len(parts) == 2 && parts[1] == normalizedRole {
				return true
			}
		}
	}
	return false
}

// HasPermission ตรวจสอบว่าผู้ใช้มี permission สำหรับ resource และ scopes ที่ระบุ
func HasPermission(c echo.Context, resourceName string, requiredScopes []string) bool {
	permissions := GetUserPermissions(c)

	for _, perm := range permissions {
		if perm.ResourceName == resourceName {
			// สร้าง scope set
			permScopeSet := make(map[string]struct{})
			for _, scope := range perm.Scopes {
				permScopeSet[strings.ToLower(strings.TrimSpace(scope))] = struct{}{}
			}

			// เช็คว่ามี scope ทุกตัวที่ต้องการ
			hasAllScopes := true
			for _, reqScope := range requiredScopes {
				normalizedScope := strings.ToLower(strings.TrimSpace(reqScope))
				if _, ok := permScopeSet[normalizedScope]; !ok {
					hasAllScopes = false
					break
				}
			}

			if hasAllScopes {
				return true
			}
		}
	}
	return false
}

// Permission middleware: ตรวจสอบว่าผู้ใช้มี permission สำหรับ resource และ scope ที่ระบุ
func (m *GoMiddleware) Permission(resourceName string, requiredScopes []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := c.Get(contextClaimsKey)
			claims, ok := val.(jwt.MapClaims)
			if !ok || claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No user claims found. Please authenticate first"})
			}

			permissions := collectPermissions(claims)
			if len(permissions) == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": fmt.Sprintf("No permissions found. Required resource: %s, scopes: %v", resourceName, requiredScopes),
				})
			}

			// หา permission ที่ตรงกับ resource name
			var matchedPermission *Permission
			for _, perm := range permissions {
				if perm.ResourceName == resourceName {
					matchedPermission = &perm
					break
				}
			}

			if matchedPermission == nil {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": fmt.Sprintf("No permission found for resource: %s", resourceName),
				})
			}

			// เช็คว่ามี scope ที่ต้องการหรือไม่
			permScopeSet := make(map[string]struct{})
			for _, scope := range matchedPermission.Scopes {
				permScopeSet[strings.ToLower(strings.TrimSpace(scope))] = struct{}{}
			}

			// ต้องมี scope ทุกตัวที่ require (all-of)
			for _, reqScope := range requiredScopes {
				normalizedScope := strings.ToLower(strings.TrimSpace(reqScope))
				if _, ok := permScopeSet[normalizedScope]; !ok {
					return c.JSON(http.StatusForbidden, map[string]string{
						"error": fmt.Sprintf("Insufficient permission scopes. Required: %v, Available: %v", requiredScopes, matchedPermission.Scopes),
					})
				}
			}

			return next(c)
		}
	}
}

// AnyPermission middleware: ตรวจสอบว่าผู้ใช้มี permission อย่างน้อย 1 อันจากที่ระบุ
func (m *GoMiddleware) AnyPermission(resourceScopeMap map[string][]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := c.Get(contextClaimsKey)
			claims, ok := val.(jwt.MapClaims)
			if !ok || claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No user claims found. Please authenticate first"})
			}

			permissions := collectPermissions(claims)
			if len(permissions) == 0 {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "No permissions found",
				})
			}

			// เช็คว่ามี permission ใดตรงตามเงื่อนไขบ้าง
			for resourceName, requiredScopes := range resourceScopeMap {
				for _, perm := range permissions {
					if perm.ResourceName == resourceName {
						// สร้าง scope set
						permScopeSet := make(map[string]struct{})
						for _, scope := range perm.Scopes {
							permScopeSet[strings.ToLower(strings.TrimSpace(scope))] = struct{}{}
						}

						// เช็คว่ามี scope ทุกตัวที่ต้องการหรือไม่
						hasAllScopes := true
						for _, reqScope := range requiredScopes {
							normalizedScope := strings.ToLower(strings.TrimSpace(reqScope))
							if _, ok := permScopeSet[normalizedScope]; !ok {
								hasAllScopes = false
								break
							}
						}

						if hasAllScopes {
							return next(c) // มี permission ที่ตรงเงื่อนไข
						}
					}
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}

// JWTCookieMiddleware: middleware ที่ตรวจสอบ JWT จาก cookies แบบ Keycloak
func (m *GoMiddleware) JWTCookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ดึง access token จาก cookie
			accessToken, err := m.cookie.GetAccessToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Access token not found",
					"message": "Please login first",
				})
			}

			// ตรวจสอบ token ด้วยวิธีเดียวกับ IsAuthorizationWithKeycloak
			// ใช้ false สำหรับ includeAudience เพื่อข้าม audience validation (comment ออกไว้ในโค้ดเดิม)
			token, err := m.parseJWTWithJWKS(c, accessToken, false)
			fmt.Println("Parsed token:", token.Valid)
			if err != nil || !token.Valid {
				// Token invalid หรือ expired - ลบ cookies
				m.cookie.ClearAllJwtCookies(c)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Invalid or expired token",
					"message": "Please login again",
				})
			}

			setTokenClaims(c, token, accessToken)

			return next(c)
		}
	}
}

// OptionalJWTCookieMiddleware: middleware ที่ไม่บังคับให้มี JWT
func (m *GoMiddleware) OptionalJWTCookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ดึง access token จาก cookie (ไม่บังคับ)
			accessToken, err := m.cookie.GetAccessToken(c)
			if err != nil {
				// ไม่มี token - ดำเนินการต่อ
				return next(c)
			}

			// ตรวจสอบ token ด้วยวิธีเดียวกับ IsAuthorizationWithKeycloak
			token, err := m.parseJWTWithJWKS(c, accessToken, true)
			if err != nil {
				// Token invalid - ล้าง cookies แล้วดำเนินการต่อ
				m.cookie.ClearAllJwtCookies(c)
				return next(c)
			}

			setTokenClaims(c, token, accessToken)

			return next(c)
		}
	}
}

// SessionMiddleware: middleware สำหรับตรวจสอบ session โดยเช็ค Session ID และ Refresh Token (Access Token เป็น optional)
func (m *GoMiddleware) SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. เช็ค session ID ก่อน (จำเป็น)
			sessionCookie, err := c.Cookie("session_id")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Session not found",
					"message": "Please login first",
				})
			}
			sessionID := sessionCookie.Value
			if sessionID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Invalid session",
					"message": "Please login first",
				})
			}

			// 2. เช็ค refresh token (จำเป็น)
			refreshToken, err := m.cookie.GetRefreshToken(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error":   "Refresh token not found",
					"message": "Please login first",
				})
			}

			// 3. เช็ค access token (ถ้ามีและ valid ให้ใช้, ถ้าไม่มีหรือหมดอายุ ก็ไม่เป็นไร)
			accessToken, err := m.cookie.GetAccessToken(c)
			if err == nil && accessToken != "" {
				// ถ้ามี access token ให้ลองตรวจสอบ (ไม่สนใจ expiry)
				token, parseErr := m.parseJWTWithoutClaimsValidation(c, accessToken)
				if parseErr != nil {
					// Token corrupted (signature ผิด) - ลบทุก cookies
					m.cookie.ClearAllJwtCookies(c)
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error":   "Invalid access token",
						"message": "Token is corrupted, please login again",
					})
				}

				// Access token valid - set claims
				setTokenClaims(c, token, accessToken)
			}
			// ถ้าไม่มี access token หรือหมดอายุ ก็ไม่เป็นไร เพราะมี refresh token อยู่

			// เก็บ session ID และ refresh token ใน context
			c.Set(contextSessionIDKey, sessionID)
			c.Set("refresh_token", refreshToken)

			return next(c)
		}
	}
}
