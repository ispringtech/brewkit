package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	defaultConfigPath = "%s/.brewkit/config"
)

var (
	DefaultConfig = Config{}
)

func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to receive user home dir")
	}

	return fmt.Sprintf(defaultConfigPath, homeDir), nil
}
