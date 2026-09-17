package main

import (
	greetv1 "connect-sample/gen/greet/v1"
	"connect-sample/gen/greet/v1/greetv1connect"
	loggerx "connect-sample/logger"
	"connect-sample/requestid"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	loggerx.Info(ctx, "request headers", slog.Any("header", req.Header()))

	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: "Hello, " + req.Msg.Name + "!",
	}), nil
}

func main() {
	logger := slog.New(loggerx.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	greeter := &GreetServer{}
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)

	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.RequestID(),
		requestid.SetRequestID(logger),
		middleware.Logger(),
	)
	e.GET("/", func(c echo.Context) error {
		loggerx.Error(c.Request().Context(), "root endpoint hit")
		return c.JSON(http.StatusOK, map[string]string{"message": "hello world"})
	})
	e.Any(fmt.Sprintf("%s*", path), echo.WrapHandler(handler))
	e.StartH2CServer("localhost:8080", &http2.Server{})
}
