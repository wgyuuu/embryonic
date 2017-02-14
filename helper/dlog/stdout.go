package dlog

import "fmt"

type Stdout struct {
	level Severity
}

var stdout *Stdout = NewStdout(STRACE)

func NewStdout(level Severity) *Stdout {
	return &Stdout{level}
}

func (s *Stdout) levelOK(level Severity) bool {
	return s.level >= level
}

func (s *Stdout) Debug(args ...interface{}) {
	s.print(DEBUG, args...)
}

func (s *Stdout) Trace(args ...interface{}) {
	s.print(TRACE, args...)
}

func (s *Stdout) Info(args ...interface{}) {
	s.print(INFO, args...)
}
func (s *Stdout) Strace(args ...interface{}) {
	s.print(STRACE, args...)
}
func (s *Stdout) Warning(args ...interface{}) {
	s.print(WARNING, args...)
}
func (s *Stdout) Error(args ...interface{}) {
	s.print(ERROR, args...)
}
func (s *Stdout) Fatal(args ...interface{}) {
	s.print(FATAL, args...)
}

func (s *Stdout) print(level Severity, args ...interface{}) {
	if !s.levelOK(level) {
		return
	}

	var typ string
	switch level {
	case FATAL:
		typ = "FATAL"
	case ERROR:
		typ = "ERROR"
	case WARNING:
		typ = "WARNING"
	case STRACE:
		typ = "STRACE"
	case INFO:
		typ = "INFO"
	case TRACE:
		typ = "TRACE"
	case DEBUG:
		typ = "DEBUG"
	}
	args = append([]interface{}{fmt.Sprintf("[%s]", typ)}, args...)
	fmt.Println(args...)
}
