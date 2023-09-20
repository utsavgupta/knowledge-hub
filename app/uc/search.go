package uc

import (
	"context"
	"fmt"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/repos"
	"github.com/utsavgupta/knowledge-hub/app/services"
)

type SearchUc func(context.Context, entities.Query) (*entities.Response, error)
type DomainStatusValidator func(context.Context, string) error

func NewSearchUc(domainStatusValidator DomainStatusValidator, responseRepo repos.ResponseRepo, conceptService services.ConceptService) SearchUc {

	return func(ctx context.Context, query entities.Query) (*entities.Response, error) {

		if err := domainStatusValidator(ctx, query.DomainId); err != nil {
			return nil, err
		}

		concepts, err := conceptService.Get(ctx, query.Question)

		if err != nil || len(concepts) < 1 {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could fetch contexts")
		}

		query.Concepts = concepts

		answer, err := responseRepo.Get(ctx, query)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			err = fmt.Errorf("could not generate answer")
		}

		return answer, err
	}
}

func NewDomainStatusValidator(resourceRepo repos.ResourceRepo) DomainStatusValidator {

	return func(ctx context.Context, domainId string) error {

		resources, err := resourceRepo.List(ctx, domainId)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return fmt.Errorf("could not fetch resources from database: %w", err)
		}

		if len(resources) < 1 {
			return fmt.Errorf("%w: either the domain id %s does not exist, or it contains no resources", ValidationError, domainId)
		}

		for _, resource := range resources {
			if resource.Status == entities.ResourceStatusIngested {
				return nil
			}
		}

		return fmt.Errorf("%w: none of the resources have been ingested for domain %s", ValidationError, domainId)
	}
}
