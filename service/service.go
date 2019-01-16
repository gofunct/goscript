package service

import (
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

var Set = wire.NewSet(
	newService,
	newOptions,
)

type Service struct {
	Pattern    string
	Version    string
	Endpoint   endpoint.Endpoint
	Middleware endpoint.Middleware
	*grpc.Server
}

func newService(pattern string, endpoint endpoint.Endpoint, middleware endpoint.Middleware, option []grpc.ServerOption) *Service {
	s := grpc.NewServer(option...)
	return &Service{Pattern: pattern, Endpoint: endpoint, Middleware: middleware, Server: s}
}

// Runnable determines if the command is itself runnable.
func (c *Service) Runnable() bool {
	return c.Pattern != "" || c.Endpoint != nil || c.Server != nil
}

func newOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	opts = append(opts, grpc.UnaryInterceptor(kitgrpc.Interceptor))

	return opts
}
