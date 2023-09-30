package logger

import (
	"fmt"
	"io"
)

type Logger interface {
	Logf(format string, a ...any)
	Outputf(format string, a ...any) // Used to output info to user
	Debugf(format string, a ...any)
}

func NewLogger(outputWriter, logWriter io.Writer, debug bool) Logger {
	return &logger{
		outputWriter: outputWriter,
		logWriter:    logWriter,
		debug:        debug,
	}
}

type logger struct {
	outputWriter io.Writer
	logWriter    io.Writer
	debug        bool
}

func (l *logger) Logf(format string, a ...any) {
	_, _ = fmt.Fprintf(l.logWriter, format, a...)
}

func (l *logger) Outputf(format string, a ...any) {
	_, _ = fmt.Fprintf(l.outputWriter, format, a...)
}

func (l *logger) Debugf(format string, a ...any) {
	if l.debug {
		_, _ = fmt.Fprintf(l.logWriter, format, a...)
	}
}
