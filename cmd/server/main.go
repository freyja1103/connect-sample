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

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	slog.InfoContext(ctx, "request headers", slog.Any("header", req.Header()))

	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: "Hello, " + req.Msg.Name + "!",
	}), nil
}

func main() {
	greeter := &GreetServer{}
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)

	e := echo.New()
	e.Any(fmt.Sprintf("%s*", path), echo.WrapHandler(handler))
	e.StartH2CServer("localhost:8080", &http2.Server{})
}
