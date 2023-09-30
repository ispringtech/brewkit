package docker

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/backend/app/docker"
	"github.com/ispringtech/brewkit/internal/common/infrastructure/executor"
	"github.com/ispringtech/brewkit/internal/common/infrastructure/logger"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/dockerfile"
)

const (
	dockerExecutable = "docker"
)

var (
	dockerEnv = map[string]string{
		"DOCKER_BUILDKIT": "1", // enable buildkit explicitly
	}
)

func NewClient(clientConfigPath maybe.Maybe[string], log logger.Logger) (docker.Client, error) {
	d, err := executor.New(
		dockerExecutable,
		executor.WithEnv(os.Environ()),
		executor.WithEnvMap(dockerEnv),
		executor.WithLogger(logger.NewExecutorLogger(log)),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		clientConfigPath: clientConfigPath,
		dockerExecutor:   d,
		outputParser:     outputParser{},
	}, nil
}

type client struct {
	clientConfigPath maybe.Maybe[string]
	dockerExecutor   executor.Executor
	outputParser     outputParser
}

func (c *client) Build(ctx context.Context, d dockerfile.Dockerfile, params docker.BuildParams) error {
	var args executor.Args

	c.populateWithCommonArgs(&args)
	c.populateWithBuilderArgs(&args)
	args.AddArgs("build")

	if maybe.Valid(params.SSHAgent) {
		args.AddKV("--ssh", fmt.Sprintf("default=%s", maybe.Just(params.SSHAgent)))
	}

	if len(params.Secrets) > 0 {
		for _, secret := range params.Secrets {
			args.AddKV("--secret", fmt.Sprintf("id=%s,src=%s", secret.ID, secret.Path))
		}
	}

	args.AddKV("--target", params.Target)

	if maybe.Valid(params.Output) {
		args.AddKV("--output", maybe.Just(params.Output))
	}

	args.AddArgs("-f-", ".") // Read Dockerfile from stdin and use PWD as context

	dockerfileReader := bytes.NewBufferString(d.Format())

	err := c.dockerExecutor.Run(ctx, args, executor.RunParams{
		Stdin: maybe.NewJust[io.Reader](dockerfileReader),
	})
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return docker.RequestError{
				Output: string(exitErr.Stderr),
				Code:   exitErr.ExitCode(),
			}
		}
	}
	return nil
}

func (c *client) Value(ctx context.Context, d dockerfile.Dockerfile, params docker.ValueParams) ([]byte, error) {
	var args executor.Args

	c.populateWithCommonArgs(&args)
	c.populateWithBuilderArgs(&args)
	args.AddArgs("build")
	args.AddKV("--progress", "plain") // Set to plain to be able to parse output

	if !params.UseCache {
		args.AddArgs("--no-cache") // Disable cache for target
	}

	if maybe.Valid(params.SSHAgent) {
		args.AddKV("--ssh", fmt.Sprintf("default=%s", maybe.Just(params.SSHAgent)))
	}

	args.AddKV("--target", params.Var)

	args.AddArgs("-f-", ".") // Read Dockerfile from stdin and use PWD as context

	dockerfileReader := bytes.NewBufferString(d.Format())

	output := &bytes.Buffer{}
	err := c.dockerExecutor.Run(ctx, args, executor.RunParams{
		Stdin:  maybe.NewJust[io.Reader](dockerfileReader),
		Stderr: maybe.NewJust[io.Writer](output),
	})
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, &docker.RequestError{
				Output: output.String(),
			}
		}
		return nil, err
	}

	return c.outputParser.parseBuildOutputForRunTarget(output)
}

func (c *client) ListImages(ctx context.Context, images []string) ([]docker.Image, error) {
	var args executor.Args

	c.populateWithCommonArgs(&args)

	args.AddArgs("image", "ls") // List images

	for _, image := range images {
		args.AddKV("--filter", fmt.Sprintf("reference=%s", image))
	}

	args.AddKV("--format", "{{.Repository}}:{{.Tag}}") // Just list images repository and tag for now, can be used as filter

	output := &bytes.Buffer{}
	err := c.dockerExecutor.Run(ctx, args, executor.RunParams{Stdout: maybe.NewJust[io.Writer](output)})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list docker images")
	}

	var res []docker.Image
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		image := scanner.Text()
		res = append(res, docker.Image{Img: image})
	}

	return res, nil
}

func (c *client) PullImage(ctx context.Context, img string) error {
	var args executor.Args

	c.populateWithCommonArgs(&args)
	args.AddArgs("pull")
	args.AddArgs(img)

	return c.dockerExecutor.Run(ctx, args, executor.RunParams{})
}

func (c *client) ClearCache(ctx context.Context, params docker.ClearCacheParams) error {
	var args executor.Args

	c.populateWithCommonArgs(&args)
	c.populateWithBuilderArgs(&args)

	args.AddArgs("prune", "-f")
	if params.All {
		args.AddArgs("-a") // Delete all cache
	}

	return c.dockerExecutor.Run(ctx, args, executor.RunParams{})
}

func (c *client) BuildImage(_ context.Context, _ string) error {
	// TODO implement me
	panic("implement me")
}

func (c *client) populateWithCommonArgs(args *executor.Args) {
	if maybe.Valid(c.clientConfigPath) {
		args.AddKV("--config", maybe.Just(c.clientConfigPath))
	}
}

func (c *client) populateWithBuilderArgs(args *executor.Args) {
	args.AddArgs("builder") // Use builder explicitly
}
