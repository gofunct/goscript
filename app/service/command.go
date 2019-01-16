package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/oklog/run"
	"net/http"
)

type Service struct {
	Name       string
	Short      string
	Long       string
	Example    string
	Version    string
	PreRun     func(ctx context.Context, s *Service)
	Endpoint   endpoint.Endpoint
	Middleware endpoint.Middleware
	Handler    http.Handler
	PostRun    func(ctx context.Context, s *Service)
	group      *run.Group
}

// Runnable determines if the command is itself runnable.
func (c *Service) Runnable() bool {
	return c.Endpoint != nil || c.Handler != nil
}
