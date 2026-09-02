package server

import (
	"log/slog"
	"os"

	mw "github.com/kassisol/tsa/api/server/middleware"
	"github.com/kassisol/tsa/api/server/router/acme"
	"github.com/kassisol/tsa/api/server/router/ca"
	"github.com/kassisol/tsa/api/server/router/crl"
	"github.com/kassisol/tsa/api/server/router/system"
	"github.com/kassisol/tsa/api/server/router/tsaca"
	"github.com/kassisol/tsa/api/server/router/usergroup"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func API(addr string, tls bool) {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer s.End()

	jwk := []byte(s.GetConfig("jwk")[0].Value)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mw.Https())
	e.Use(mw.AdminPassword())
	e.Use(mw.CAInit())

	// Directory
	e.GET("/", system.IndexHandle)

	// Authz
	h := middleware.BasicAuth(mw.Authorization)(system.AuthzHandle)
	e.GET("/new-authz", h)

	// Version
	e.GET("/version", system.ServerVersionHandle)

	// System
	jwtConfig := echojwt.Config{
		Skipper:    mw.DefaultSkipper,
		SigningKey: jwk,
	}

	sys := e.Group("/system")
	sys.Use(echojwt.WithConfig(jwtConfig))
	sys.Use(mw.AdminOnly())

	sys.GET("/info", system.InfoHandle)
	sys.PUT("/admin/password", system.AdminChangePasswordHandle)
	sys.POST("/ca/init", system.CAInitHandle)

	sys.GET("/auth", system.AuthListHandle)
	sys.POST("/auth", system.AuthCreateHandle)
	sys.DELETE("/auth/:key", system.AuthDeleteHandle)
	sys.PUT("/auth/enable/:type", system.AuthEnableHandle)
	sys.PUT("/auth/disable", system.AuthDisableHandle)

	sys.GET("/cert", system.CertListHandle)
	sys.DELETE("/cert/revoke/:serialnumber", system.CertRevokeHandle)
	sys.DELETE("/cert/revoke-by-cn/:cn", system.CertRevokeByCNHandle)

	// CA public certificate
	e.GET("/ca", ca.PubCertHandle)

	// Revocation file
	e.GET("/crl/CRL.crl", crl.CRLHandle)

	// ACME
	r := e.Group("/acme")
	r.Use(echojwt.JWT(jwk))

	// New certificate
	r.POST("/new-app", acme.NewCertHandle)

	// Revoke
	r.POST("/revoke-cert", acme.RevokeCertHandle)

	// User Groups
	ug := e.Group("/tsa/ug")
	ug.Use(echojwt.JWT(jwk))
	ug.Use(mw.AdminOnly())

	ug.POST("/register", usergroup.RegisterHandle)
	ug.DELETE("/register", usergroup.UnregisterHandle)

	// TSA CA
	tc := e.Group("/tsa/ca")
	tc.Use(echojwt.JWT(jwk))
	tc.Use(mw.AdminOnly())

	tc.DELETE("/cert", tsaca.CertDeleteHandle)

	if tls {
		if err := e.StartTLS(addr, cfg.API.CrtFile, cfg.API.KeyFile); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	} else {
		if err := e.Start(addr); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
}
