package middleware

import (
	"context"

	helperCookie "github.com/GodeFvt/go-backend/helper/cookie"
	helperRedis "github.com/GodeFvt/go-backend/redis"
	"github.com/labstack/echo/v4"
)

type GoMiddlewareInf interface {
	InitContextIfNotExists(next echo.HandlerFunc) echo.HandlerFunc
	SetTracer(next echo.HandlerFunc) echo.HandlerFunc

	IsAuthorizationWithKeycloak(next echo.HandlerFunc) echo.HandlerFunc
	IsAuthorizationBasicJwt(next echo.HandlerFunc) echo.HandlerFunc

	JWTCookieMiddleware() echo.MiddlewareFunc
	OptionalJWTCookieMiddleware() echo.MiddlewareFunc

	SessionMiddleware() echo.MiddlewareFunc

	Role(requiredRoles []string) echo.MiddlewareFunc
	Permission(resourceName string, requiredScopes []string) echo.MiddlewareFunc

	ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc

	ValidateParamId(key string) echo.MiddlewareFunc
}

type GoMiddleware struct {
	ctx         context.Context
	redis       *helperRedis.Client
	cookie      *helperCookie.CookieManager
	jwtSecret   string
	jwksURL     []string
	allowedAlgs []string
	audience    []string
	issuer      string
}

func InitMiddleware(redis *helperRedis.Client, cookie *helperCookie.CookieManager, jwtSecret string, jwksURL []string, allowedAlgs []string, audience []string, issuer string) GoMiddlewareInf {
	return &GoMiddleware{
		ctx:         context.TODO(),
		redis:       redis,
		cookie:      cookie,
		jwtSecret:   jwtSecret,
		jwksURL:     jwksURL,
		allowedAlgs: allowedAlgs,
		audience:    audience,
		issuer:      issuer,
	}
}

func (m *GoMiddleware) InitContextIfNotExists(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		if ctx == nil {
			bgCtx := context.Background()
			newReq := c.Request().WithContext(bgCtx)

			c.SetRequest(newReq)
		}
		return next(c)
	}
}
