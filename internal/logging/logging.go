package logging

import (
	"fmt"
	"log/slog"
)

type Logger struct {
	Service string
}

func NewLogger(service string) *Logger {
	return &Logger{
		Service: service,
	}
}

func (l *Logger) Start() {
	fmt.Println()
	slog.Info("##### " + l.Service + " #####")
	fmt.Println()
}

func (l *Logger) Finish() {
	slog.Info(l.Service + " Finished!")
}
