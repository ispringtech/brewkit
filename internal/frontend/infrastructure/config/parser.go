package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/google/go-jsonnet"
	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/common/slices"
	"github.com/ispringtech/brewkit/internal/frontend/app/config"
)

type Parser struct{}

func (p Parser) Config(configPath string) (config.Config, error) {
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config.Config{}, errors.WithStack(config.ErrConfigNotFound)
		}

		return config.Config{}, errors.Wrap(err, "failed to read config file")
	}

	vm := jsonnet.MakeVM()

	jsonnet.Version()

	data, err := vm.EvaluateAnonymousSnippet(path.Base(configPath), string(fileBytes))
	if err != nil {
		return config.Config{}, errors.Wrap(err, "failed to compile jsonnet for config")
	}

	var c Config

	err = json.Unmarshal([]byte(data), &c)
	if err != nil {
		return config.Config{}, errors.Wrap(err, "failed to parse json config")
	}

	return config.Config{
		Secrets: slices.Map(c.Secrets, func(s Secret) config.Secret {
			return config.Secret{
				ID:   s.ID,
				Path: os.ExpandEnv(s.Path),
			}
		}),
	}, nil
}

func (p Parser) Dump(srcConfig config.Config) ([]byte, error) {
	c := Config{
		Secrets: slices.Map(srcConfig.Secrets, func(s config.Secret) Secret {
			return Secret{
				ID:   s.ID,
				Path: s.Path,
			}
		}),
	}

	data, err := json.Marshal(c)
	return data, errors.Wrap(err, "failed to marshal config to json")
}
