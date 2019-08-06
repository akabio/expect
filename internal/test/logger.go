package test

import (
	"fmt"
	"testing"

	"gitlab.com/testle/expect"
)

// Logger implements Test and allows for inspection of the
// calls.
type Logger struct {
	Fatal    []string
	Error    []string
	Messages []string
	t        *testing.T
}

func New(t *testing.T) *Logger {
	return &Logger{t: t}
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

// ExpectMessages returns the messages as a expect-value
func (l *Logger) ExpectMessages() expect.Val {
	return expect.Value(l.t, "messages", l.Messages)
}

// ExpectMessage returns the message at given index
func (l *Logger) ExpectMessage(i int) expect.Val {
	if len(l.Messages) <= i {
		l.t.Errorf("there is not message at index %v", i)
	}
	return expect.Value(l.t, "message", l.Messages[i])
}
