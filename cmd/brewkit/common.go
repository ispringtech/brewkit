package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ispringtech/brewkit/internal/common/infrastructure/logger"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	appconfig "github.com/ispringtech/brewkit/internal/frontend/app/config"
	infraconfig "github.com/ispringtech/brewkit/internal/frontend/infrastructure/config"
)

type commonOpt struct {
	configPath             string
	verbose                bool
	dockerClientConfigPath maybe.Maybe[string]
}

func (o *commonOpt) scan(ctx *cli.Context) {
	o.configPath = ctx.String("config")
	o.verbose = ctx.Bool("verbose")
	dockerConfigPath := ctx.String("docker-config")
	if dockerConfigPath != "" {
		o.dockerClientConfigPath = maybe.NewJust(dockerConfigPath)
	}
}

func makeLogger(verbose bool) logger.Logger {
	return logger.NewLogger(os.Stdout, os.Stderr, verbose)
}

func parseConfig(configPath string, log logger.Logger) (appconfig.Config, error) {
	c, err := infraconfig.Parser{}.Config(configPath)
	if err != nil {
		if !errors.Is(err, appconfig.ErrConfigNotFound) {
			return appconfig.Config{}, err
		}
		log.Debugf("config not found in %s: default config will be used\n", configPath)

		return appconfig.DefaultConfig, nil
	}
	return c, err
}
