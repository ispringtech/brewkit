package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ispringtech/brewkit/internal/frontend/infrastructure/jsonnet"
)

func fmtCommand() *cli.Command {
	return &cli.Command{
		Name:   "fmt",
		Usage:  "jsonnetfmt passed files",
		Action: executeFmt,
	}
}

func executeFmt(ctx *cli.Context) error {
	formatter := jsonnet.Formatter{}

	for _, filepath := range ctx.Args().Slice() {
		fileInfo, err := os.Stat(filepath)
		if err != nil {
			return err
		}

		format, err := formatter.Format(filepath)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath, []byte(format), fileInfo.Mode())
		if err != nil {
			return err
		}
	}

	return nil
}
