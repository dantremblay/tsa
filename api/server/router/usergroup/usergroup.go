package usergroup

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/juliengk/stack/jsonapi"
	"github.com/kassisol/tsa/api/auth"
	"github.com/kassisol/tsa/api/storage"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/tsa/pkg/api"
	"github.com/labstack/echo/v4"
)

func RegisterHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	// Read plain string body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot read request body")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	group := strings.TrimSpace(string(body))
	if group == "" {
		r := jsonapi.NewErrorResponse(1000, "Group value cannot be empty")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	// Verify auth backend is enabled
	authTypeConfig := s.GetConfig("auth_type")
	if len(authTypeConfig) == 0 {
		r := jsonapi.NewErrorResponse(1000, "Make sure to enable a backend auth before adding configuration")
		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	authType := authTypeConfig[0].Value
	if authType == "none" {
		r := jsonapi.NewErrorResponse(1000, "Make sure to enable a backend auth before adding configuration")
		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	a, err := auth.NewDriver(authType)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	if err = a.AddConfig("auth_group_user", group); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusUnprocessableEntity, r)
	}

	response := jsonapi.NewSuccessResponseWithMessage(nil, 1001, fmt.Sprintf("User group registered: %s", group))
	return api.JSON(c, http.StatusCreated, response)
}

func UnregisterHandle(c echo.Context) error {
	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, err.Error())
		return api.JSON(c, http.StatusInternalServerError, r)
	}
	defer s.End()

	// Read plain string body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		r := jsonapi.NewErrorResponse(1000, "Cannot read request body")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	group := strings.TrimSpace(string(body))
	if group == "" {
		r := jsonapi.NewErrorResponse(1000, "Group value cannot be empty")
		return api.JSON(c, http.StatusBadRequest, r)
	}

	s.RemoveConfig("auth_group_user", group)

	return c.NoContent(http.StatusNoContent)
}
