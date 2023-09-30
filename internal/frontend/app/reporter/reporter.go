package reporter

type Reporter interface {
	Logf(format string, a ...any)
	Debugf(format string, a ...any)
}
