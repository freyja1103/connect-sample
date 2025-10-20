package main

import (
	greetv1 "connect-sample/gen/greet/v1"
	"connect-sample/gen/greet/v1/greetv1connect"
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/http2"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, in *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	req := new(connect.Request[greetv1.GreetRequest])
	slog.InfoContext(ctx, "request headers", slog.Any("header", req.Header()))

	return &greetv1.GreetResponse{
		Greeting: "Hello, " + in.Name + "!",
	}, nil
}

func main() {
	greeter := &GreetServer{}
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)

	e := echo.New()
	e.Any(fmt.Sprintf("%s*", path), echo.WrapHandler(handler))
	e.StartH2CServer("localhost:8080", &http2.Server{})
}
