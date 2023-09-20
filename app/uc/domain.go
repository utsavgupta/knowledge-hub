package uc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/repos"
)

var (
	domainIdRegEx = regexp.MustCompile("^[A-Z][a-z]{,4}([a-z]+?_*?[a-z]+?){,10}?")
)

type ListDomainsUc func(context.Context) ([]entities.Domain, error)
type AddDomainUc func(context.Context, entities.Domain) (*entities.Domain, error)
type DeleteDomainUc func(context.Context, string) error

func NewListDomainsUc(repo repos.DomainRepo) ListDomainsUc {

	return func(ctx context.Context) ([]entities.Domain, error) {

		entities, err := repo.List(ctx)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could not fetch domain list")
		}

		return entities, nil
	}
}

func NewAddDomainUc(repo repos.DomainRepo) AddDomainUc {

	return func(ctx context.Context, domain entities.Domain) (*entities.Domain, error) {

		if err := validateDomainEntity(domain); err != nil {
			logger.Instance().Debug(ctx, err.Error())
			return nil, err
		}

		if ent, _ := repo.Get(ctx, domain.Id); ent != nil {
			return nil, fmt.Errorf("%w: domain id already exists", ValidationError)
		}

		domain.CreatedAt = time.Now()

		ent, err := repo.Create(ctx, domain)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return nil, fmt.Errorf("could not create domain")
		}

		return ent, nil
	}
}

func NewDeleteDomainUc(repo repos.DomainRepo) DeleteDomainUc {

	return func(ctx context.Context, id string) error {

		err := repo.Delete(ctx, id)

		if err != nil {
			logger.Instance().Error(ctx, err.Error())
			return fmt.Errorf("could not delete domain")
		}

		return err
	}
}

func validateDomainEntity(domain entities.Domain) error {

	if !domainIdRegEx.MatchString(domain.Id) || len(domain.Id) > 15 {
		return fmt.Errorf("%w: the id can be at most 15 characters long. it should start will an upper case character and may contain an underscore.", ValidationError)
	}

	if len(domain.Name) > 50 {
		return fmt.Errorf("%w: the name can be 50 characters long.", ValidationError)
	}

	if len(domain.Name) > 140 {
		return fmt.Errorf("%w: the description can be 140 characters long.", ValidationError)
	}

	return nil
}
