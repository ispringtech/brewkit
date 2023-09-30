package logger

import (
	"github.com/ispringtech/brewkit/internal/common/infrastructure/executor"
)

func NewExecutorLogger(logger Logger) executor.Logger {
	return &executorLogger{logger: logger}
}

type executorLogger struct {
	logger Logger
}

func (logger *executorLogger) Info(s string) {
	logger.logger.Logf(s)
}

func (logger *executorLogger) Debug(s string) {
	logger.logger.Debugf(s)
}
