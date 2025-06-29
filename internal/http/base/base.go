package base

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ResponseRoot struct {
	Total *uint64          `json:"total,omitempty"`
	Data  *json.RawMessage `json:"data"`
}

type ErrorResponse struct {
	Details []*ErrorDetail `json:"details,omitempty"`
	Status  int            `json:"status"`
	Title   string         `json:"title"`
}

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Text  string `json:"text"`
}

func JSON(c echo.Context, code int, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}

	obj := json.RawMessage(b)

	return c.JSON(code, &ResponseRoot{Data: &obj})
}

func JSONWithTotal(c echo.Context, code int, data interface{}, total uint64) error {
	b, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}

	obj := json.RawMessage(b)

	return c.JSON(code, &ResponseRoot{Data: &obj, Total: &total})
}

func Bind(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		log.Error().Stack().Err(errors.WithStack(err)).Msg("")

		code := http.StatusBadRequest

		if err := c.JSON(code, &ErrorResponse{
			Status:  code,
			Title:   http.StatusText(code),
			Details: []*ErrorDetail{{Text: "invalid parameter"}},
		}); err != nil {
			log.Error().Stack().Err(errors.WithStack(err)).Msg("")
		}

		return err
	}

	return nil
}

// HandleError handles domain errors and returns appropriate HTTP responses
func HandleError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	code := http.StatusInternalServerError
	details := []*ErrorDetail{{Text: err.Error()}}

	// Handle validation errors
	switch err.Error() {
	case "name is required":
		code = http.StatusBadRequest
		details = []*ErrorDetail{{Field: "name", Text: err.Error()}}
	case "email is required":
		code = http.StatusBadRequest
		details = []*ErrorDetail{{Field: "email", Text: err.Error()}}
	case "password is required":
		code = http.StatusBadRequest
		details = []*ErrorDetail{{Field: "password", Text: err.Error()}}
	case "password must be at least 8 characters long":
		code = http.StatusBadRequest
		details = []*ErrorDetail{{Field: "password", Text: err.Error()}}
	default:
		// Handle UUID parsing errors
		if strings.Contains(err.Error(), "invalid UUID") {
			code = http.StatusBadRequest
			details = []*ErrorDetail{{Text: "invalid UUID format"}}
		}
		// Handle not found errors
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			code = http.StatusNotFound
			details = []*ErrorDetail{{Text: "User not found"}}
		}
		// Handle duplicate key errors
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint") {
			code = http.StatusConflict
			details = []*ErrorDetail{{Text: "Resource already exists"}}
		}
	}

	if err := c.JSON(code, &ErrorResponse{
		Status:  code,
		Title:   http.StatusText(code),
		Details: details,
	}); err != nil {
		log.Error().Stack().Err(errors.WithStack(err)).Msg("")
	}

	return err
}
