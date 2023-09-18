package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/utsavgupta/knowledge-hub/app/adapters/datasources"
	"github.com/utsavgupta/knowledge-hub/app/adapters/transport"
	"github.com/utsavgupta/knowledge-hub/app/repos"
	"github.com/utsavgupta/knowledge-hub/app/runners"
	"github.com/utsavgupta/knowledge-hub/app/services"
	"github.com/utsavgupta/knowledge-hub/app/uc"
)

func configureHttpRunner() (runners.Runner, error) {

	var postgresConnString string
	var weaviateHost *url.URL
	var openaiAccessKey string
	var port int
	var err error

	if postgresConnString, err = getStringFromEnv("kh_pg_conn_str"); err != nil {
		return nil, err
	}

	if weaviateHost, err = getURLFromEnv("kh_weaviate_host"); err != nil {
		return nil, err
	}

	if openaiAccessKey, err = getStringFromEnv("kh_openai_access_key"); err != nil {
		return nil, err
	}

	if port, err = getIntFromEnv("kh_app_port"); err != nil {
		return nil, err
	}

	runnerDependencies, err := createHttpRunnerDependencies(postgresConnString, weaviateHost, openaiAccessKey)

	if err != nil {
		return nil, err
	}

	return transport.NewHttpRunner(port, *runnerDependencies), nil
}

func createHttpRunnerDependencies(postgresConnString string, weaviateHost *url.URL, openaiAccessKey string) (*transport.HttpRunnerDependencies, error) {

	var err error
	var pgConnPool *pgxpool.Pool

	if pgConnPool, err = createPgConnectionPool(postgresConnString); err != nil {
		return nil, err
	}

	var domainRepo repos.DomainRepo
	var resourceRepo repos.ResourceRepo
	var answerRepo repos.AnswerRepo
	var conceptService services.ConceptService

	conceptService = datasources.NewConceptOpenAI(http.DefaultClient, openaiAccessKey)

	if domainRepo, err = datasources.NewPGDomainRepo(pgConnPool); err != nil {
		return nil, err
	}

	if resourceRepo, err = datasources.NewPGResourceRepo(pgConnPool); err != nil {
		return nil, err
	}

	if answerRepo, err = datasources.NewWeaviateAnswerRepo(weaviateHost.Scheme, weaviateHost.Host, openaiAccessKey); err != nil {
		return nil, err
	}

	return &transport.HttpRunnerDependencies{
		SearchUc:         uc.NewSearchUc(domainRepo, answerRepo, conceptService),
		ListDomainsUc:    uc.NewListDomainsUc(domainRepo),
		AddDomainUc:      uc.NewAddDomainUc(domainRepo),
		DeleteDomainUc:   uc.NewDeleteDomainUc(domainRepo),
		ListResourcesUc:  uc.NewListResourcesUc(resourceRepo),
		AddResourceUc:    uc.NewAddResourceUc(resourceRepo, domainRepo),
		DeleteResourceUc: uc.NewDeleteResourceUc(resourceRepo),
	}, nil
}

func createPgConnectionPool(connStr string) (*pgxpool.Pool, error) {

	return pgxpool.New(context.Background(), connStr)
}

func getStringFromEnv(name string) (string, error) {

	if v, ok := os.LookupEnv(name); ok {

		return v, nil
	}

	return "", fmt.Errorf("environment variable %s not found", name)
}

func getIntFromEnv(name string) (int, error) {

	if v, ok := os.LookupEnv(name); ok {

		if num, err := strconv.Atoi(v); err == nil {

			return num, nil
		}

		return -1, fmt.Errorf("environment variable %s is not an integer", name)
	}

	return -1, fmt.Errorf("environment variable %s not found", name)
}

func getURLFromEnv(name string) (*url.URL, error) {

	if v, ok := os.LookupEnv(name); ok {

		if uri, err := url.Parse(v); err == nil {

			return uri, nil
		}

		return nil, fmt.Errorf("environment variable %s is not a url", name)
	}

	return nil, fmt.Errorf("environment variable %s not found", name)
}
