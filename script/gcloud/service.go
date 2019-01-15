package gcloud

import (
	"bytes"
	"context"
	"github.com/gofunct/goscript/script"
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	osexec "os/exec"
	"os/signal"
	"sync"
)

type Gcloud struct {
	Project string
	*iio.Service
}

// New returns a new Interface which will os/exec to run commands.
// New returns a new Interface which will os/exec to run commands.
func New(i *iio.Service) *Gcloud {
	return &Gcloud{
		Service: nil,
	}
}

var Set = wire.NewSet(
	New,
	iio.Provider,
	wire.Bind(new(script.Executor), new(Gcloud)),
)

func (e *Gcloud) Exec(ctx context.Context, opts ...script.Option) (out []byte, err error) {
	var wg sync.WaitGroup

	c := script.BuildCommand("gcloud", opts)

	cmd := osexec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = c.Env

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer recover()
		for sig := range sigCh {
			log.Println("signal received--*", sig)
			if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
				break
			}
			cmd.Process.Signal(sig)
		}
	}()

	out, err = e.exec(c, cmd)
	if err != nil {
		err = errors.WithStack(err)
	}

	signal.Reset()
	close(sigCh)

	wg.Wait()
	return
}

func (g *Gcloud) exec(c *script.Command, cmd *osexec.Cmd) (out []byte, err error) {
	if c.IOConnected {
		var (
			buf bytes.Buffer
			wg  sync.WaitGroup
		)

		closers := make([]func() error, 0, 2)

		outReader, eerr := cmd.StdoutPipe()
		if eerr != nil {
			err = errors.WithStack(eerr)
			return
		}
		errReader, eerr := cmd.StderrPipe()
		if eerr != nil {
			err = errors.WithStack(eerr)
			return
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			io.Copy(g.Out(), io.TeeReader(outReader, &buf))
		}()
		closers = append(closers, outReader.Close)
		go func() {
			defer wg.Done()
			io.Copy(g.Err(), io.TeeReader(errReader, &buf))
		}()
		closers = append(closers, errReader.Close)

		cmd.Stdin = g.In()

		err = cmd.Run()
		for _, c := range closers {
			c()
		}
		wg.Wait()

		out = buf.Bytes()
	} else {
		out, err = cmd.CombinedOutput()
	}

	return
}

func (d *Gcloud) Run(ctx context.Context, opts ...script.Option) error {
	var wg sync.WaitGroup

	c := script.BuildCommand("gcloud", opts)

	cmd := osexec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = c.Env

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer recover()
		for sig := range sigCh {
			log.Println("signal received--*", sig)
			if cmd.ProcessState == nil || cmd.ProcessState.Exited() {
				break
			}
			cmd.Process.Signal(sig)
		}
	}()

	err := d.run(c, cmd)
	if err != nil {
		err = errors.WithStack(err)
	}

	signal.Reset()
	close(sigCh)

	wg.Wait()
	return nil
}

func (d *Gcloud) run(c *script.Command, cmd *osexec.Cmd) error {
	var (
		buf bytes.Buffer
		wg  sync.WaitGroup
	)

	closers := make([]func() error, 0, 2)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithStack(err)

	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return errors.WithStack(err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(d.Out(), io.TeeReader(outReader, &buf))
	}()
	closers = append(closers, outReader.Close)
	go func() {
		defer wg.Done()
		io.Copy(d.Err(), io.TeeReader(errReader, &buf))
	}()
	closers = append(closers, errReader.Close)

	cmd.Stdin = d.In()

	return cmd.Run()
}
