package apis

import (
	"log"
	"net/http"
	"strings"

	"qvault/services/secret"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		userID, err := secret.Authenticate(strings.TrimSpace(parts[1]))
		if err == secret.ErrInvalidAPIKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}
		if err != nil {
			log.Printf("authenticate request failed: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		c.Set("user_id", userID)

		return next(c)
	}
}
