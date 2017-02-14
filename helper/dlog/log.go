package dlog

type Severity int

const (
	FATAL Severity = iota
	ERROR
	WARNING
	STRACE
	INFO
	TRACE
	DEBUG
)

var logger Logger

type Logger interface {
	Debug(args ...interface{})
	Trace(args ...interface{})
	Info(args ...interface{})
	Strace(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

func getLogger() Logger {
	if logger == nil {
		return stdout
	}
	return logger
}

func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}
func Trace(args ...interface{}) {
	getLogger().Trace(args...)
}
func Info(args ...interface{}) {
	getLogger().Info(args...)
}
func Strace(args ...interface{}) {
	getLogger().Strace(args...)
}
func Warning(args ...interface{}) {
	getLogger().Warning(args...)
}
func Error(args ...interface{}) {
	getLogger().Error(args...)
}
func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}
