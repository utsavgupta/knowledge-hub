package services

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type ConceptService interface {
	Get(context.Context, string) ([]entities.Concept, error)
}
