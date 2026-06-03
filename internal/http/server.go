package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/SoraDaibu/go-clean-starter/builder"
	imiddleware "github.com/SoraDaibu/go-clean-starter/internal/http/middleware"
)

type Server struct {
	closer   func() error
	echo     *echo.Echo
	port     uint16
	logLevel string
}

func NewServer(d *builder.Dependency) *Server {
	s := &Server{
		port:     d.Config.App.ListenPort,
		logLevel: d.Config.App.LogLevel,
	}

	s.closer = func() error {
		d.DB.Close()
		return nil
	}

	s.echo = setup(d)

	return s
}

func (s *Server) Close() error {
	return s.closer()
}

func (s *Server) Run() {
	if s.logLevel == "debug" {
		//nolint:errchkjson
		data, _ := json.MarshalIndent(s.echo.Routes(), "", "  ")
		fmt.Println(string(data))
	}

	s.echo.Logger.Fatal(s.echo.Start(fmt.Sprintf(":%d", s.port)))
}

func setup(d *builder.Dependency) *echo.Echo {
	e := echo.New()

	d.Logger.Info("logger configured", "level", d.Config.App.LogLevel)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(
		imiddleware.Recover(),
		middleware.Logger(),
		middleware.RequestID(),
		// RequestLogger must run after RequestID so the request ID is available.
		imiddleware.RequestLogger(d.Logger),
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
