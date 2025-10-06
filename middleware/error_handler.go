package middleware

import (
	"errors"
	"net/http"

	"github.com/IT-CP25-US1-School-Management-System/sms-data-service/errs"
	"github.com/labstack/echo/v4"
)

func (m *GoMiddleware) ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			return c.JSON(appErr.HTTPStatus, map[string]interface{}{
				"code":    appErr.HTTPStatus,
				"message": appErr.Message,
			})
		}

		if he, ok := err.(*echo.HTTPError); ok {
			return c.JSON(he.Code, map[string]interface{}{
				"code":    he.Code,
				"message": he.Message,
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
}
