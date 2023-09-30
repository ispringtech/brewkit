package executor

import (
	"fmt"
)

type options struct {
	env    []string
	logger Logger
}

type Opt interface {
	apply(*options)
}

type optFunc func(*options)

func (f optFunc) apply(o *options) {
	f(o)
}

func WithEnv(env []string) Opt {
	return optFunc(func(o *options) {
		o.env = append(o.env, env...)
	})
}

func WithEnvMap(envMap EnvMap) Opt {
	return optFunc(func(o *options) {
		o.env = append(o.env, envMap.slice()...)
	})
}

func WithLogger(logger Logger) Opt {
	return optFunc(func(o *options) {
		o.logger = logger
	})
}

type EnvMap map[string]string

func (e EnvMap) slice() []string {
	res := make([]string, 0, len(e))
	for k, v := range e {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res
}
