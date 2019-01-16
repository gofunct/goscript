package script

import (
	"github.com/fatih/color"
	"io"
	"os"
)

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (c *Command) SetOutput(output io.Writer) {
	c.output = output
}

// OutOrStdout returns output to stdout.
func (c *Command) OutOrStdout() io.Writer {
	return c.getOut(os.Stdout)
}

// OutOrStderr returns output to stderr
func (c *Command) OutOrStderr() io.Writer {
	return c.getOut(os.Stderr)
}

func (c *Command) getOut(def io.Writer) io.Writer {
	if c.output != nil {
		return c.output
	}
	if c.HasParent() {
		return c.parent.getOut(def)
	}
	return def
}

// Println is a convenience method to Println to the defined output, fallback to Stderr if not set.
func (c *Command) Debug(msg string, args ...string) {
	color.BlueString(msg, args)
}
