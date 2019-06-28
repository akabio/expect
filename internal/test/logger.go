package test

import "fmt"

// Logger implements Test and allows for inspection of the
// calls.
type Logger struct {
	Fatal    []string
	Error    []string
	Messages []string
}

// Fatalf records call
func (l *Logger) Fatalf(f string, i ...interface{}) {
	line := fmt.Sprintf(f, i...)
	l.Fatal = append(l.Fatal, line)
	l.Messages = append(l.Messages, line)
}

// Errorf records call
func (l *Logger) Errorf(f string, i ...interface{}) {
	line := fmt.Sprintf(f, i...)
	l.Error = append(l.Error, line)
	l.Messages = append(l.Messages, line)
}
