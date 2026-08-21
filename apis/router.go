package apis

import (
	"fmt"
	"log"
	"os"

	"qvault/apis/health"
	"qvault/apis/secret"
	"qvault/services/common"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Run() error {
	if _, err := common.GetMasterKey(); err != nil {
		return fmt.Errorf("invalid MASTER_KEY: %w", err)
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit(1 * 1024 * 1024))

	api := e.Group("/api/v1")
	api.GET("/health", health.Health)

	secretGroup := api.Group("/secrets", AuthMiddleware)
	secretGroup.GET("", secret.ListKeys)
	secretGroup.GET("/:key", secret.GetKey)
	secretGroup.PUT("/:key", secret.PutKey)
	secretGroup.DELETE("/:key", secret.DeleteKey)

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	listenAddr := host + ":" + port
	log.Println("Listening on", listenAddr)

	return e.Start(listenAddr)
}
