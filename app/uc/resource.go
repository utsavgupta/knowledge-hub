package uc

import (
	"context"
	"fmt"
	"time"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/repos"
)

type ListResourcesUc func(context.Context, string) ([]entities.Resource, error)
type AddResourceUc func(context.Context, entities.Resource) (*entities.Resource, error)
type DeleteResourceUc func(context.Context, int) error

func NewListResourcesUc(repo repos.ResourceRepo) ListResourcesUc {

	return func(ctx context.Context, domainId string) ([]entities.Resource, error) {

		entities, err := repo.List(ctx, domainId)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could not fetch resources list")
		}

		return entities, nil
	}
}

func NewAddResourceUc(resourceRepo repos.ResourceRepo, domainRepo repos.DomainRepo) AddResourceUc {

	return func(ctx context.Context, resource entities.Resource) (*entities.Resource, error) {

		if err := validateResourceEntity(resource); err != nil {
			return nil, err
		}

		if ent, _ := domainRepo.Get(ctx, resource.DomainId); ent == nil {
			return nil, fmt.Errorf("%w: invalid domain id", ValidationError)
		}

		resource.CreatedAt = time.Now()

		ent, err := resourceRepo.Create(ctx, resource)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could not create resource")
		}

		return ent, nil
	}
}

func NewDeleteResourceUc(repo repos.ResourceRepo) DeleteResourceUc {

	return func(ctx context.Context, id int) error {

		err := repo.Delete(ctx, id)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return fmt.Errorf("could not delete resource")
		}

		return err
	}
}

func validateResourceEntity(resource entities.Resource) error {

	if len(resource.DomainId) < 2 || len(resource.DomainId) > 15 {
		return fmt.Errorf("%w: domain id can beetween 2 and 15 characters long.", ValidationError)
	}

	if len(resource.Name) <= 50 {
		return fmt.Errorf("%w: the name can be 50 characters long", ValidationError)
	}

	if len(resource.Description) <= 140 {
		fmt.Errorf("%w: the description can be 140 characters long", ValidationError)
	}

	return nil
}
