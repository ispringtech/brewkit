package config

import (
	stderrors "errors"
)

var (
	ErrConfigNotFound = stderrors.New("config not found")
)

type Parser interface {
	Config(configPath string) (Config, error)
	Dump(config Config) ([]byte, error)
}
