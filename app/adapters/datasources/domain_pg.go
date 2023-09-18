package datasources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/repos"
)

type pgDomainRepo struct {
	conn *pgxpool.Conn
}

func NewPGDomainRepo(connPool *pgxpool.Conn) (repos.DomainRepo, error) {

	return &pgDomainRepo{connPool}, nil
}

func (repo *pgDomainRepo) List(ctx context.Context) ([]entities.Domain, error) {

	var domains []entities.Domain

	row, err := repo.conn.Query(ctx, "SELECT id, name, description, created_at, updated_at FROM domains")

	if err != nil {

		return nil, fmt.Errorf("could not list domains: %w", err)
	}

	defer row.Close()

	domains = make([]entities.Domain, 0)

	for row.Next() {

		domain := entities.Domain{}

		if err = row.Scan(&domain.Id, &domain.Name, &domain.Description, &domain.CreatedAt, &domain.UpdatedAt); err != nil {

			return nil, fmt.Errorf("could not read domain: %w", err)
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

func (repo *pgDomainRepo) Get(ctx context.Context, id string) (*entities.Domain, error) {

	var domain entities.Domain

	row, err := repo.conn.Query(ctx, "SELECT id, name, description, created_at, updated_at FROM domains WHERE id = $1", id)

	if err != nil {

		return nil, fmt.Errorf("could not fetch domain with id %s: %w", id, err)
	}

	defer row.Close()

	if !row.Next() {

		return nil, nil
	}

	if err = row.Scan(&domain.Id, &domain.Name, &domain.Description, &domain.CreatedAt, &domain.UpdatedAt); err != nil {

		return nil, fmt.Errorf("could not fetch task with id %s: %w", id, err)
	}

	return &domain, nil
}

func (repo *pgDomainRepo) Create(ctx context.Context, domain entities.Domain) (*entities.Domain, error) {

	_, err := repo.conn.Exec(ctx, "INSERT INTO domains (id, name, description, created_at) VALUES ($1, $2, $3, $4)", domain.Id, domain.Name, domain.Description, domain.CreatedAt)

	if err != nil {

		err = fmt.Errorf("could not create domain %v: %w", domain, err)
	}

	return &domain, nil
}

func (repo *pgDomainRepo) Update(ctx context.Context, domain entities.Domain) (*entities.Domain, error) {

	_, err := repo.conn.Exec(ctx, "UPDATE domains SET name = $2, description = $3, updated_at = $4 WHERE id = $1", domain.Id, domain.Name, domain.Description, domain.UpdatedAt)

	if err != nil {

		err = fmt.Errorf("could not update domain %v: %w", domain, err)
	}

	return &domain, err
}

func (repo *pgDomainRepo) Delete(ctx context.Context, id string) error {

	var err error

	if _, err = repo.conn.Exec(ctx, "DELETE FROM domains WHERE id = $1", id); err != nil {

		err = fmt.Errorf("could not delete domain with id %s: %w", id, err)
	}

	return err
}
