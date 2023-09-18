package repos

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type DomainRepo interface {
	List(context.Context) ([]entities.Domain, error)
	Get(context.Context, string) (*entities.Domain, error)
	Create(context.Context, entities.Domain) (*entities.Domain, error)
	Update(context.Context, entities.Domain) (*entities.Domain, error)
	Delete(context.Context, string) error
}
