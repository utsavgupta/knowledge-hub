package logger

import (
	"context"
	"sync"
)

var (
	instance  Logger
	onceMutex sync.Once
)

type Logger interface {
	Info(context.Context, string)
	Debug(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
	SetLevel(string)
}

func Instance() Logger {

	return instance
}

func InitLogger(logger Logger) {

	onceMutex.Do(func() {
		instance = logger
	})
}
