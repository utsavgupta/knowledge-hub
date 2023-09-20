package repos

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type ResponseRepo interface {
	Get(context.Context, entities.Query) (*entities.Response, error)
}
