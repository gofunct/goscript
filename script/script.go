package script

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Executor is an interface for executing external commands.
type Executor interface {
	Exec(ctx context.Context, opts ...Option) ([]byte, error)
	Run(ctx context.Context, opts ...Option) error
}

// Command contains parameters for executing external commands.
type Command struct {
	Name        string
	Args        []string
	Dir         string
	Env         []string
	IOConnected bool
}

// Option specifies external command execution configurations.
type Option func(*Command)

// BuildCommand create a new Command object from given options.
func BuildCommand(name string, opts []Option) *Command {
	c := &Command{Name: name, Env: os.Environ()}
	for _, f := range opts {
		f(c)
	}
	return c
}

// WithArgs sets arguments for a command.
func WithArgs(args ...string) Option {
	return func(c *Command) {
		c.Args = append(c.Args, args...)
	}
}

// WithDir sets a working directory for a command.
func WithDir(dir string) Option {
	return func(c *Command) {
		c.Dir = dir
	}
}

// WithPATH sets a PATH for a command.
func WithPATH(value string) Option {
	return func(c *Command) {
		for i := 0; i < len(c.Env); i++ {
			kv := strings.Split(c.Env[i], "=")
			if kv[0] == "PATH" {
				c.Env[i] = "PATH=" + value + string(filepath.ListSeparator) + kv[1]
				return
			}
		}
		WithEnv("PATH", value)(c)
	}
}

// WithEnv append a environment variable for a command.
func WithEnv(key, value string) Option {
	return func(c *Command) {
		c.Env = append(c.Env, key+"="+value)
	}
}

// WithIOConnected connects i/o with a command.
func WithIOConnected() Option {
	return func(c *Command) {
		c.IOConnected = true
	}
}

func Handle(executor Executor, opts ...Option) Handler {
	return func(ctx context.Context, i interface{}) (interface{}, error) {
		return executor.Exec(ctx, opts...)
	}
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
