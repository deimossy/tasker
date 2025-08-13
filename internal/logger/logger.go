package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	ch chan string
}

func NewLogger(bufferSize int) *Logger {
	l := &Logger{
		ch: make(chan string, bufferSize),
	}

	go l.run()
	return l
}

func (l *Logger) run() {
	for msg := range l.ch {
		fmt.Println(time.Now().Format(time.RFC3339), msg)
	}
}

func (l *Logger) Log(msg string) {
	select {
	case l.ch <- msg:
	default:
	}
}

func (l *Logger) Close() {
	close(l.ch)
}