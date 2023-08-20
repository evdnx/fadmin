package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/evdnx/unixmint/auth"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// check if Authorization header exists in the request
		authHeader := c.Get("Authorization", "")
		if authHeader == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// decode token brancaToken
		_, err := auth.DecodeToken(authHeader, 24, nil)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}
