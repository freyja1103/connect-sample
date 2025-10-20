package main

import (
	greetv1 "connect-sample/gen/greet/v1"
	"connect-sample/gen/greet/v1/greetv1connect"
	"context"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	serveMux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	serveMux.Handle(path, handler)
	http.ListenAndServe("localhost:8080", h2c.NewHandler(serveMux, &http2.Server{}))
}
