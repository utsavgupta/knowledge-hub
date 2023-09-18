package entities

import (
	"time"
)

const (
	ResourceStatusNew       = "NEW"
	ResourceStatusIngesting = "INGESTING"
	ResourceStatusIngested  = "INGESTED"
)

type Resource struct {
	Id                   int        `json:"id"`
	DomainId             string     `json:"domainId"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Status               string     `json:"status"`
	Url                  string     `json:"url"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt,omitempty"`
	IngestionStartedAt   *time.Time `json:"ingestion_started_at,omitempty"`
	IngestionCompletedAt *time.Time `json:"ingestion_completed_at,omitempty"`
}
