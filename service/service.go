package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

var Set = wire.NewSet(
	newService,
)

type Service struct {
	Pattern    string
	Version    string
	Endpoint   endpoint.Endpoint
	Middleware endpoint.Middleware
	Handler    http.Handler
}

func newService(pattern string, endpoint endpoint.Endpoint, middleware endpoint.Middleware, handler http.Handler) *Service {
	return &Service{Pattern: pattern, Endpoint: endpoint, Middleware: middleware, Handler: handler}
}

// Runnable determines if the command is itself runnable.
func (c *Service) Runnable() bool {
	return c.Pattern != "" || c.Endpoint != nil || c.Handler != nil
}

func (s *Service) HandleGrpc(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			s.Handler.ServeHTTP(w, r)
		}
	})
}
