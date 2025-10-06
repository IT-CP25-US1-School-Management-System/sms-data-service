package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddL "github.com/labstack/echo/v4/middleware"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	AllowCredentials bool
	MaxAge           int
}

func SetupCORS(config CORSConfig) echo.MiddlewareFunc {
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
			"X-Auth",
		}
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		}
	}

	if config.MaxAge == 0 {
		config.MaxAge = 86400
	}

	return echoMiddL.CORSWithConfig(echoMiddL.CORSConfig{
		Skipper:          echoMiddL.DefaultSkipper,
		AllowOrigins:     config.AllowOrigins,
		AllowHeaders:     config.AllowHeaders,
		AllowMethods:     config.AllowMethods,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	})
}

func DefaultCORSConfig(allowOrigins, allowOriginHeaders []string, allowCredentials bool) echo.MiddlewareFunc {
	defaultAllowOriginHeaders := []string{
		echo.HeaderOrigin,
		echo.HeaderContentType,
		echo.HeaderAccept,
		echo.HeaderAuthorization,
		"X-Requested-With",
		"X-Auth",
	}
	if len(allowOriginHeaders) > 0 {
		defaultAllowOriginHeaders = append(defaultAllowOriginHeaders, allowOriginHeaders...)
	}

	return SetupCORS(CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     defaultAllowOriginHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           86400,
	})
}
