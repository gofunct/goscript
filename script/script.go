package script

import (
	"context"
	"encoding/json"
	"net/http"
)

// Executor is an interface for executing external commands.
type Executor interface {
	Exec(ctx context.Context) ([]byte, error)
	Run(ctx context.Context) error
}

func Handle(executor Executor) Handler {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return executor.Exec(ctx)
	}
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
