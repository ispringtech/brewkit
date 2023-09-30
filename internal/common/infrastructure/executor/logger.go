package executor

type Logger interface {
	Info(s string)
	Debug(s string)
}
