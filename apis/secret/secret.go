package secret

import (
	"net/http"
	"strings"

	"qvault/services/secret"

	"github.com/labstack/echo/v5"
)

func ListKeys(c *echo.Context) error {
	keys, err := secret.ListKeys(userID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, keys)
}

func GetKey(c *echo.Context) error {
	key, err := parseKey(c)
	if err != nil {
		return err
	}

	sec, err := secret.GetSecret(userID(c), key)
	if err == secret.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err == secret.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "secret not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sec)
}

func PutKey(c *echo.Context) error {
	key, err := parseKey(c)
	if err != nil {
		return err
	}

	var request struct {
		Value string `json:"value"`
	}
	if err := c.Bind(&request); err != nil || strings.TrimSpace(request.Value) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	sec, err := secret.PutSecret(userID(c), key, request.Value)
	if err == secret.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err == secret.ErrSecretValueTooLarge {
		return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "secret value is too large")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sec)
}

func DeleteKey(c *echo.Context) error {
	key, err := parseKey(c)
	if err != nil {
		return err
	}

	err = secret.DeleteSecret(userID(c), key)
	if err == secret.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err == secret.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "secret not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
