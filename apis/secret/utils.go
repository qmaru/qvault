package secret

import (
	"net/http"
	"net/url"
	"strings"

	"qkms/services/kms"

	"github.com/labstack/echo/v5"
)

func parseKey(c *echo.Context) (string, error) {
	key, err := url.PathUnescape(strings.TrimSpace(c.Param("key")))
	if err != nil || kms.ValidateSecretKey(key) != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "secret key is required")
	}
	return key, nil
}

func userID(c *echo.Context) int64 {
	value := c.Get("user_id")
	if id, ok := value.(int64); ok {
		return id
	}
	return 0
}
