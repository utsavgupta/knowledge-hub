package uc

import (
	"context"

	"github.com/utsavgupta/knowledge-hub/app/entities"
)

type SearchUc func(context.Context, entities.Query) (*entities.Answer, error)

func NewSearchUc() SearchUc {

	return func(ctx context.Context, query entities.Query) (*entities.Answer, error) {

		return nil, nil
	}
}
