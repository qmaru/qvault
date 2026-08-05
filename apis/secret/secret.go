package secret

import (
	"net/http"
	"strings"

	"qkms/services/kms"

	"github.com/labstack/echo/v5"
)

func ListKeys(c *echo.Context) error {
	keys, err := kms.ListKeys(userID(c))
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

	secret, err := kms.GetSecret(userID(c), key)
	if err == kms.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err == kms.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "secret not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, secret)
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

	secret, err := kms.PutSecret(userID(c), key, request.Value)
	if err == kms.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, secret)
}

func DeleteKey(c *echo.Context) error {
	key, err := parseKey(c)
	if err != nil {
		return err
	}

	err = kms.DeleteSecret(userID(c), key)
	if err == kms.ErrInvalidSecretKey {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid secret key")
	}
	if err == kms.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, "secret not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
