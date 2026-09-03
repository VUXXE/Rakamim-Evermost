package e2e

import (
	"net/http"
	"testing"
	"time"

	"evermos-api/internal/infrastructure/container"
	"evermos-api/internal/infrastructure/mysql"
	httpserver "evermos-api/internal/server/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *fiber.App {
	cfg := mysql.Config{
		Host:            "127.0.0.1",
		Port:            "3306",
		User:            "root",
		Password:        "password",
		DBName:          "evermos",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := mysql.Connect(cfg)
	require.NoError(t, err)

	err = mysql.AutoMigrate(db)
	require.NoError(t, err)

	c := container.SetupContainer(db)

	app := fiber.New()
	httpserver.SetupRoutes(httpserver.RouterConfig{
		App:       app,
		Container: c,
	})

	return app
}

func testReq(app *fiber.App, req *http.Request) (*http.Response, error) {
	return app.Test(req, 10000)
}
