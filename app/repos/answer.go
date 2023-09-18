package repos

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type AnswerRepo interface {
	Get(context.Context, entities.Query) (*entities.Answer, error)
}
