package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	helperCookie "github.com/GodeFvt/go-backend/helper/cookie"
	helperRedis "github.com/GodeFvt/go-backend/redis"
	"github.com/joncalhoun/qson"
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

	InputForm(next echo.HandlerFunc) echo.HandlerFunc
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

func (m *GoMiddleware) InputForm(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := Form(c); err != nil {
			var code int
			var message interface{}
			if he, ok := err.(*echo.HTTPError); ok {
				code = he.Code
				message = he.Message
			}
			return echo.NewHTTPError(code, message)
		}
		return next(c)
	}
}

func Form(c echo.Context) error {
	var data = map[string]interface{}{}
	reqMethod := c.Request().Method
	Header := c.Request().Header

	if reqMethod == http.MethodPost || reqMethod == http.MethodPut || reqMethod == http.MethodDelete {
		contentType := Header.Get("Content-Type")
		fmt.Println("Content-Type:", contentType)
		if strings.Contains(contentType, echo.MIMEMultipartForm) {
			fmt.Println("Parsing multipart/form-data")
			form, err := c.MultipartForm()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": http.ErrMissingBoundary.Error() + " or has not any parameter"})
			}
			bu, _ := qson.ToJSON(url.Values(form.Value).Encode())
			json.Unmarshal(bu, &data)
			fmt.Println("Form values parsed:", data)

			data, err = parseOnKeyData(data)
			fmt.Println("Parsed data from multipart/form-data:", data)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			}

			for k, v := range form.File {
				c.Set(k, v)
			}

		} else if strings.Contains(strings.ToLower(contentType), echo.MIMEApplicationJSON) {
			var err error
			if err := json.NewDecoder(c.Request().Body).Decode(&data); err != nil && err != io.EOF {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			}
			data, err = parseOnKeyData(data)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			}
		} else if strings.Contains(strings.ToLower(contentType), echo.MIMEApplicationForm) {
			postForm, err := c.FormParams()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			}
			if reqMethod == http.MethodDelete {
				buf := bytes.Buffer{}
				io.Copy(&buf, c.Request().Body)
				postForm, _ = url.ParseQuery(buf.String())
			}
			if len(postForm) > 0 {
				bu, _ := qson.ToJSON(postForm.Encode())
				json.Unmarshal(bu, &data)
			}
			data, err = parseOnKeyData(data)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{"message": err.Error()})
			}
		}
	}

	if len(data) > 0 {
		fmt.Println("Parsed form data:", data)
		c.Set("params", data)
	}
	return nil
}

func parseOnKeyData(data map[string]interface{}) (map[string]interface{}, error) {
	if data != nil {
		/*
			support on data from json format
		*/
		if v, ok := data["data"]; ok {
			valueType := reflect.ValueOf(v).Kind()
			switch valueType {
			case reflect.Map:
				data = v.(map[string]interface{})
			case reflect.String:
				data = map[string]interface{}{}
				if err := json.Unmarshal([]byte(v.(string)), &data); err != nil {
					return data, err
				}
			}
		}
	}

	return data, nil
}
