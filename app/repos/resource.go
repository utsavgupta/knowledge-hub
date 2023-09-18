package repos

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type ResourceRepo interface {
	List(context.Context, string) ([]entities.Resource, error)
	Get(context.Context, int) (*entities.Resource, error)
	Create(context.Context, entities.Resource) (*entities.Resource, error)
	Update(context.Context, entities.Resource) (*entities.Resource, error)
	Delete(context.Context, int) error
}
