package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/SoraDaibu/go-clean-starter/internal/http/base"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

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
					log.Error().Stack().Err(err).Msg("")
				}

				const code = http.StatusInternalServerError

				err := c.JSON(code, &base.ErrorResponse{
					Status: code,
					Title:  http.StatusText(code),
				})
				if err != nil {
					log.Error().Err(err).Msg("")
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
			log.Debug().Str("request_body", string(reqBody)).Msg("Request body")
		} else {
			log.Debug().Msg("Request: Binary")
		}

		log.Debug().Str("response_body", string(resBody)).Msg("Response body")
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
