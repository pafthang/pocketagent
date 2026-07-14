package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

type httpError struct {
	code    int
	message string
}

func (e *httpError) Error() string {
	return e.message
}

// ErrBadRequest returns a 400-class handler error.
func ErrBadRequest(msg string) error { return &httpError{code: 400, message: msg} }

// ErrForbidden returns a 403-class handler error.
func ErrForbidden(msg string) error { return &httpError{code: 403, message: msg} }

// ErrNotFound returns a 404-class handler error.
func ErrNotFound(msg string) error { return &httpError{code: 404, message: msg} }

// ErrConflict returns a 409-class handler error.
func ErrConflict(msg string) error { return &httpError{code: 409, message: msg} }

// ErrInternal returns a 500-class handler error.
func ErrInternal(msg string) error { return &httpError{code: 500, message: msg} }

// WrapInternal wraps an error as 500-class.
func WrapInternal(action string, err error) error {
	return ErrInternal(fmt.Sprintf("%s: %v", action, err))
}

// RespondError maps domain errors to JSON HTTP responses.
func RespondError(c echo.Context, err error) error {
	var he *httpError
	if errors.As(err, &he) {
		return c.JSON(he.code, map[string]string{"error": he.message})
	}
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

// MapPocketError maps PocketBase client errors to HTTP responses.
func MapPocketError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	var apiErr *pbclient.APIError
	if errors.As(err, &apiErr) {
		status := apiErr.StatusCode
		if status < 400 {
			status = http.StatusInternalServerError
		}
		return c.JSON(status, map[string]string{"error": apiErr.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

// MapMemoError maps memo client errors to HTTP responses.
func MapMemoError(c echo.Context, err error) error {
	msg := err.Error()
	if strings.Contains(msg, "not found") {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "document not found"})
	}
	return c.JSON(http.StatusBadGateway, map[string]string{"error": msg})
}