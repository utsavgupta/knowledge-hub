package datasources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/repos"
)

type pgResourceRepo struct {
	conn *pgxpool.Conn
}

func NewPGResourceRepo(connPool *pgxpool.Conn) (repos.ResourceRepo, error) {

	return &pgResourceRepo{connPool}, nil
}

func (repo *pgResourceRepo) List(ctx context.Context, domainId string) ([]entities.Resource, error) {

	var resources []entities.Resource

	row, err := repo.conn.Query(ctx, "SELECT id, name, description, status, url,domain_id, created_at, updated_at, ingestion_started_at, ingestion_completed_at FROM resources WHERE domain_id = $1", domainId)

	if err != nil {

		return nil, fmt.Errorf("could not list resources: %w", err)
	}

	defer row.Close()

	resources = make([]entities.Resource, 0)

	for row.Next() {

		resource := entities.Resource{}

		if err = row.Scan(&resource.Id, &resource.Name, &resource.Description, &resource.Status, &resource.Url, &resource.DomainId, &resource.CreatedAt, &resource.UpdatedAt, &resource.IngestionStartedAt, &resource.IngestionCompletedAt); err != nil {

			return nil, fmt.Errorf("could not read resource: %w", err)
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (repo *pgResourceRepo) Get(ctx context.Context, id int) (*entities.Resource, error) {

	var resource entities.Resource

	row, err := repo.conn.Query(ctx, "SELECT id, name, description, status, url,domain_id, created_at, updated_at, ingestion_started_at, ingestion_completed_at FROM resources where id = $1", id)

	if err != nil {

		return nil, fmt.Errorf("could not fetch resource with id %d: %w", id, err)
	}

	defer row.Close()

	if !row.Next() {

		return nil, nil
	}

	if err = row.Scan(&resource.Id, &resource.Name, &resource.Description, &resource.Status, &resource.Url, &resource.DomainId, &resource.CreatedAt, &resource.UpdatedAt, &resource.IngestionStartedAt, &resource.IngestionCompletedAt); err != nil {

		return nil, fmt.Errorf("could not fetch resource with id %d: %w", id, err)
	}

	return &resource, nil
}

func (repo *pgResourceRepo) Create(ctx context.Context, resource entities.Resource) (*entities.Resource, error) {

	_, err := repo.conn.Exec(ctx, "INSERT INTO resources (name, description, status, url,domain_id, created_at) VALUES ($1, $2, $3, $4)", resource.Name, resource.Description, entities.ResourceStatusNew, resource.Url.RequestURI(), resource.DomainId, resource.CreatedAt)

	if err != nil {

		err = fmt.Errorf("could not create resource %v: %w", resource, err)
	}

	return &resource, nil
}

func (repo *pgResourceRepo) Update(ctx context.Context, resource entities.Resource) (*entities.Resource, error) {

	_, err := repo.conn.Exec(ctx,
		"UPDATE domains SET name = $2, description = $3 , updated_at = $4 WHERE id = $1",
		resource.Id, resource.Name, resource.Description, resource.UpdatedAt)

	if err != nil {

		err = fmt.Errorf("could not update resource %v: %w", resource, err)
	}

	return &resource, err
}

func (repo *pgResourceRepo) Delete(ctx context.Context, id int) error {

	var err error

	if _, err = repo.conn.Exec(ctx, "DELETE FROM resources WHERE id = $1", id); err != nil {

		err = fmt.Errorf("could not delete resource with id %d: %w", id, err)
	}

	return err
}
