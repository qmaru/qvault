package apis

import (
	"net/http"
	"strings"

	"qkms/services/kms"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		userID, err := kms.Authenticate(strings.TrimSpace(parts[1]))
		if err == kms.ErrInvalidAPIKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Set("user_id", userID)

		return next(c)
	}
}
