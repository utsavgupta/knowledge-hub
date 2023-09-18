package main

import (
	"context"
	"os"

	"github.com/utsavgupta/knowledge-hub/app/logger"
)

func main() {

	logger.InitLogger(logger.NewZeroLogger())

	runner, err := configureHttpRunner()

	if err != nil {
		logger.Instance().Error(context.Background(), err.Error())
		os.Exit(1)
	}

	runner.Run()
}
