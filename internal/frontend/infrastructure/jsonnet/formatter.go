package jsonnet

import (
	"os"
	"path"

	jsonnetformatter "github.com/google/go-jsonnet/formatter"
	"github.com/pkg/errors"
)

type Formatter struct{}

func (formatter Formatter) Format(configPath string) (string, error) {
	filename := path.Base(configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	options := jsonnetformatter.DefaultOptions()
	options.Indent = 4
	options.StringStyle = jsonnetformatter.StringStyleLeave

	data, err := jsonnetformatter.Format(filename, string(configData), options)
	return data, errors.Wrap(err, "failed to format config file")
}
