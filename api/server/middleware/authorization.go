package middleware

import (
	"log/slog"
	"os"

	pass "github.com/juliengk/go-utils/password"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/auth/driver"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/labstack/echo/v4"
)

func Authorization(username, password string, c echo.Context) (bool, error) {
	var loginStatus driver.LoginStatus

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return false, err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer s.End()

	authType := s.GetConfig("auth_type")[0].Value
	if authType == "none" {
		slog.Warn("No authentication configured")
	}

	if username == "admin" {
		if pass.ComparePassword([]byte(password), []byte(s.GetConfig("admin_password")[0].Value)) {
			loginStatus = 1
		}
	} else {
		a, err := auth.NewDriver(authType)
		if err != nil {
			slog.Warn(err.Error())
		}

		loginStatus, err = a.Login(username, password)
		if err != nil {
			slog.Warn(err.Error())

			return false, err
		}
	}

	if loginStatus > 0 {
		c.Set("username", username)

		admin := false
		if loginStatus == 1 {
			admin = true
		}
		c.Set("admin", admin)

		return true, nil
	}

	return false, nil
}
