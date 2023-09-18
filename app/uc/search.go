package uc

import (
	"context"
	"fmt"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/repos"
	"github.com/utsavgupta/knowledge-hub/app/services"
)

type SearchUc func(context.Context, entities.Query) (*entities.Answer, error)

func NewSearchUc(domainRepo repos.DomainRepo, answerRepo repos.AnswerRepo, conceptService services.ConceptService) SearchUc {

	return func(ctx context.Context, query entities.Query) (*entities.Answer, error) {

		domain, err := domainRepo.Get(ctx, query.DomainId)

		if err != nil {
			return nil, fmt.Errorf("could not fetch domain from database: %w", err)
		}

		if domain == nil {
			return nil, fmt.Errorf("%w: invalid domain id", ValidationError)
		}

		concepts, err := conceptService.Get(ctx, query.Question)

		if err != nil || len(concepts) < 1 {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could fetch contexts")
		}

		query.Concepts = concepts

		answer, err := answerRepo.Get(ctx, query)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			err = fmt.Errorf("could not generate answer")
		}

		return answer, err
	}
}
