package script

import (
	"context"
	context2 "golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
type Handler func(ctx context.Context, request interface{}) (response interface{}, err error)

// Nop is an endpoint that does nothing and returns a nil error.
// Useful for tests.
func Nop(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil }

// Middleware is a chainable behavior modifier for endpoints.
type Middleware func(Handler) Handler

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

type ScriptHandler struct {}

func NewScriptHandler() *ScriptHandler {
	return &ScriptHandler{}
}

func (s *ScriptHandler) Exec(ctx context2.Context, cmd *Command) (*Output, error) {
	e := exec.CommandContext(ctx, cmd.Name, cmd.Args...)
	log.Print("starting command")
	cmd.Env = os.Environ()
	data, err := e.Output()
	log.Println(string(data))
	return &Output{
		Out:                  data,
	}, err
}

