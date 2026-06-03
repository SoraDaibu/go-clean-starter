package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/SoraDaibu/go-clean-starter/internal/http/base"
	"github.com/SoraDaibu/go-clean-starter/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestLogger derives a per-request logger from baseLogger, tagged with the request
// ID set by echo's RequestID middleware, and stores it in the request context so
// downstream HTTP call sites can retrieve it via logger.FromContext.
func RequestLogger(baseLogger logger.Logger) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Response().Header().Get(echo.HeaderXRequestID)
			req := c.Request()
			ctx := logger.IntoContext(req.Context(), baseLogger.With("request_id", id))
			c.SetRequest(req.WithContext(ctx))

			return h(c)
		}
	}
}

func Recover() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}

				msgs := []string{}

				for depth := 0; ; depth++ {
					_, file, line, ok := runtime.Caller(depth)
					if !ok {
						break
					}

					msgs = append(msgs, fmt.Sprintf(
						"======> %d: %v:%d", depth, file, line,
					))
				}

				errs := []error{}
				if e, ok := rec.(error); ok {
					errs = append(errs, e)
				} else {
					errs = append(errs, fmt.Errorf("%+v", rec))
				}

				errs = append(errs, errors.New(strings.Join(msgs, "\n")))

				for _, err := range errs {
					logger.FromContext(c.Request().Context()).Error("panic recovered", "err", err)
				}

				const code = http.StatusInternalServerError

				err := c.JSON(code, &base.ErrorResponse{
					Status: code,
					Title:  http.StatusText(code),
				})
				if err != nil {
					logger.FromContext(c.Request().Context()).Error("failed to write recover response", "err", err)
				}
			}()

			return h(c)
		}
	}
}

func BodyDump(env string) echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if env == "production" {
			return
		}

		if c.Request().Header.Get(echo.HeaderContentType) == "application/json" {
			logger.FromContext(c.Request().Context()).Debug("Request body", "request_body", string(reqBody))
		} else {
			logger.FromContext(c.Request().Context()).Debug("Request: Binary")
		}

		logger.FromContext(c.Request().Context()).Debug("Response body", "response_body", string(resBody))
	})
}

func DefaultContentType() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get(echo.HeaderContentType) == "" {
				c.Request().Header.Set(echo.HeaderContentType, "application/json")
			}

			return h(c)
		}
	}
}
