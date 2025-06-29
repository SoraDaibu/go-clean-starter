package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/SoraDaibu/go-clean-starter/builder"
	imiddleware "github.com/SoraDaibu/go-clean-starter/internal/http/middleware"
)

type Server struct {
	closer func() error
	echo   *echo.Echo
	port   uint16
}

func NewServer(d *builder.Dependency) *Server {
	s := &Server{port: d.Config.App.ListenPort}

	s.closer = func() error {
		return d.DB.Close()
	}

	s.echo = setup(d)

	return s
}

func (s *Server) Close() error {
	return s.closer()
}

func (s *Server) Run() {
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		//nolint:errchkjson
		data, _ := json.MarshalIndent(s.echo.Routes(), "", "  ")
		fmt.Println(string(data))
	}

	s.echo.Logger.Fatal(s.echo.Start(fmt.Sprintf(":%d", s.port)))
}

func setup(d *builder.Dependency) *echo.Echo {
	e := echo.New()

	var level zerolog.Level
	switch d.Config.App.LogLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warning", "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// To show file:line where log was called
	zerolog.CallerSkipFrameCount = 2

	// output in JSON
	writer := os.Stdout

	// Create structured logger with timestamp and file:line
	log.Logger = zerolog.New(writer).
		With().
		Timestamp(). // Add ISO timestamp
		Caller().    // Show file:line where log was called
		Logger()

	log.Info().Str("level", level.String()).Msg("Zerolog configured")

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(
		imiddleware.Recover(),
		middleware.Logger(),
		middleware.RequestID(),
		middleware.Secure(),
		imiddleware.DefaultContentType(),
		imiddleware.BodyDump(d.Config.App.Env),
	)

	registerRoutes(d, e)

	return e
}

func registerRoutes(d *builder.Dependency, e *echo.Echo) {
	// health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	{
		// users
		user := e.Group("/users")
		userHandler := builder.InitializeUserHandler(d)

		user.GET("/:id", userHandler.GetUser)
		user.POST("", userHandler.CreateUser)
	}
}
