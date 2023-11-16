package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/evdnx/unixmint/auth"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check if Authorization header exists in the request
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			// decode token brancaToken
			_, err := auth.DecodeToken(authHeader, 24, nil)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
