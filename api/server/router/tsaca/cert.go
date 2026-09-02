package tsaca

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ocsp"
)

func CertDeleteHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	// Read plain string body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot read request body")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	cn := strings.TrimSpace(string(body))
	if cn == "" {
		r := jsonapi.NewErrorResponse(1000, "CN value cannot be empty")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	certs := db.List(map[string]string{"cn": cn, "status": "V"})

	if len(certs) == 0 {
		r := jsonapi.NewErrorResponse(1000, "No valid certificate found for CN: "+cn)
		return api.JSON(c, http.StatusNotFound, r)
	}

	revocationDate := ca.DatabaseDateTimeFormat(time.Now())
	revocationReason := ocsp.CessationOfOperation

	for _, cert := range certs {
		db.Revoke(cert.SerialNumber, revocationDate, revocationReason)
	}

	return c.NoContent(http.StatusNoContent)
}
