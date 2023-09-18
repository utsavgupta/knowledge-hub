package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/utsavgupta/knowledge-hub/app/adapters/datasources"
	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
)

func main() {

	logger.InitLogger(logger.NewZeroLogger())

	conceptService := datasources.NewConceptOpenAI(http.DefaultClient, "")
	answerRepo, err := datasources.NewWeaviateAnswerRepo("http", "localhost:8080", "")

	if err != nil {
		panic(err)
	}

	question := "What are the tax benefits offered under section 54?"

	concepts, err := conceptService.Get(context.Background(), question)

	if err != nil {
		logger.Instance().Error(context.Background(), err.Error())
		os.Exit(1)
	}

	logger.Instance().Info(context.Background(), fmt.Sprintf("Concepts: %v", concepts))

	query := entities.Query{
		Question: question,
		DomainId: "Cleartax_agent",
		Concepts: concepts,
	}

	_, err = answerRepo.Get(context.Background(), query)

	if err != nil {
		logger.Instance().Error(context.Background(), err.Error())
		os.Exit(1)
	}
}
