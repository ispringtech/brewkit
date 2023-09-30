package executor

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/common/maybe"
)

type RunParams struct {
	Stdin  maybe.Maybe[io.Reader]
	Stdout maybe.Maybe[io.Writer]
	Stderr maybe.Maybe[io.Writer]
}

type Executor interface {
	Run(ctx context.Context, args Args, params RunParams) error
}

func New(executable string, opts ...Opt) (Executor, error) {
	_, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}

	o := options{}
	for _, opt := range opts {
		opt.apply(&o)
	}

	return &executor{
		executable: executable,
		options:    o,
	}, nil
}

type executor struct {
	executable string
	options    options
}

func (e *executor) Run(ctx context.Context, args Args, params RunParams) (err error) {
	cmd := exec.Command(e.executable, args...) // #nosec G204

	cmd.Stdin = maybe.MapNone(params.Stdin, func() io.Reader {
		return os.Stdin
	})
	cmd.Stdout = maybe.MapNone(params.Stdout, func() io.Writer {
		return os.Stdout
	})
	cmd.Stderr = maybe.MapNone(params.Stderr, func() io.Writer {
		return os.Stderr
	})

	cmd.Env = e.options.env

	err = e.logArgs(args)
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "failed to run cmd %s", cmd.String())
	}

	commandCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var runErr error
	go func() {
		runErr = cmd.Wait()
		cancelFunc()
	}()

	select {
	case <-ctx.Done():
		if e.options.logger != nil {
			e.options.logger.Info("Stopping process: docker\n")
		}

		killErr := cmd.Process.Signal(os.Interrupt)
		err = errors.WithStack(ctx.Err())
		if killErr != nil {
			err = errors.Wrapf(err, killErr.Error())
		}
		return err
	case <-commandCtx.Done():
		err = runErr
		return err
	}
}

func (e *executor) logArgs(args []string) error {
	if e.options.logger == nil {
		return nil
	}

	bytes, err := json.Marshal(args)
	if err != nil {
		return errors.Wrap(err, "failed to marshal command args")
	}

	e.options.logger.Debug(string(bytes) + "\n")

	return nil
}
