package acme

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/juliengk/go-cert/ca"
	"github.com/juliengk/go-cert/ca/database"
	"github.com/juliengk/go-cert/errors"
	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/server/httputils"
	"github.com/kassisol/tsa/api/types"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/kassisol/tsa/pkg/token"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ocsp"
)

func RevokeCertHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	db, err := database.NewBackend("sqlite", cfg.CA.Dir.Root)
	if err != nil {
		e := errors.New(errors.CertStoreError, errors.ReadFailed)
		r := jsonapi.NewErrorResponse(e.ErrorCode, e.Message)

		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer db.End()

	// Get JWT Claims
	authHeader := c.Request().Header.Get("Authorization")
	jwt, _ := token.JWTFromHeader(authHeader, "Bearer")

	jwk, err := httputils.GetTokenSigningKey()
	if err != nil {
		return api.JSON(c, http.StatusInternalServerError, err)
	}

	t := token.New(jwk, true)
	claims, _ := t.GetCustomClaims(jwt)

	// Get POST data
	revokecert := new(types.RevokeCert)

	if err := c.Bind(revokecert); err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot unmarshal JSON")

		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	// Validate
	rcert := db.List(map[string]string{"serial": strconv.Itoa(revokecert.SerialNumber)})[0]

	if rcert.StatusFlag != "V" {
		var msg string
		switch rcert.StatusFlag {
		case "R":
			msg = "Certificate already revoked"
		case "E":
			msg = "Certificate has expired"
		default:
			msg = "Invalid revocation status"
		}
		r := jsonapi.NewErrorResponse(int(errors.OCSPError)+int(errors.InvalidStatus), msg)

		return api.JSON(c, http.StatusBadRequest, r)
	}

	reCN := regexp.MustCompile(`CN=([a-z0-9\.\-\_]+)$`)

	cn := reCN.FindStringSubmatch(rcert.DistinguishedName)[1]

	if cn != claims.GetFirstAudience() && !claims.Admin {
		r := jsonapi.NewErrorResponse(11000, "Cannot revoke a certificate for which you are not the owner")

		return api.JSON(c, http.StatusBadRequest, r)
	}

	// Revoke certificate
	revocationDate := ca.DatabaseDateTimeFormat(time.Now())
	revocationReason := ocsp.CessationOfOperation

	db.Revoke(revokecert.SerialNumber, revocationDate, revocationReason)

	// Response
	response := jsonapi.NewSuccessResponseWithMessage(nil, 1001, "Certificate revoked")

	return api.JSON(c, http.StatusOK, response)
}
