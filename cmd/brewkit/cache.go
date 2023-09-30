package main

import (
	"github.com/urfave/cli/v2"

	backendcache "github.com/ispringtech/brewkit/internal/backend/app/cache"
	"github.com/ispringtech/brewkit/internal/backend/infrastructure/docker"
	"github.com/ispringtech/brewkit/internal/frontend/app/service"
)

func cache() *cli.Command {
	return &cli.Command{
		Name:  "cache",
		Usage: "Manipulate brewkit docker cache",
		Subcommands: []*cli.Command{
			cacheClear(),
		},
	}
}

func cacheClear() *cli.Command {
	return &cli.Command{
		Name:  "clear",
		Usage: "Clear docker builder cache",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "Clear all cache, not just dangling ones",
			},
		},
		Action: func(ctx *cli.Context) error {
			var opts commonOpt
			opts.scan(ctx)
			clearAll := ctx.Bool("all")

			logger := makeLogger(opts.verbose)

			dockerClient, err := docker.NewClient(opts.dockerClientConfigPath, logger)
			if err != nil {
				return err
			}

			cacheAPI := backendcache.NewCacheService(dockerClient)
			cacheService := service.NewCacheService(cacheAPI)

			return cacheService.ClearCache(ctx.Context, service.ClearCacheParam{
				All: clearAll,
			})
		},
	}
}
