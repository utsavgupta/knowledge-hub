package datasources

import (
	"context"
	"fmt"
	"strings"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/repos"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	graphqlModels "github.com/weaviate/weaviate/entities/models"
)

type weaviateResponseRepo struct {
	client *weaviate.Client
}

func NewWeaviateResponseRepo(scheme string, host string, openaiAccessKey string) (repos.ResponseRepo, error) {

	cfg := weaviate.Config{
		Host:    host,
		Scheme:  scheme,
		Headers: map[string]string{"X-OpenAI-Api-Key": openaiAccessKey},
	}

	client, err := weaviate.NewClient(cfg)

	if err != nil {
		return nil, fmt.Errorf("could not create Weaviate client with config %v: %w", cfg, err)
	}

	return &weaviateResponseRepo{client}, nil
}

func (repo *weaviateResponseRepo) Get(ctx context.Context, query entities.Query) (*entities.Response, error) {

	generativeSearchBuilder := graphql.NewGenerativeSearch().GroupedResult(query.Question)
	nearTextArgumentBuilder := repo.prepareNearTextArgumentBuilder(query.Concepts)

	gqlResponse, err := repo.client.GraphQL().
		Get().
		WithClassName(query.DomainId).
		WithFields(graphql.Field{Name: "source"}).
		WithGenerativeSearch(generativeSearchBuilder).
		WithNearText(nearTextArgumentBuilder).
		WithLimit(5).Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve answer from Weaviate for question `%s`: %w", query.Question, err)
	}

	return repo.prepareResponseFromGQLResponse(query, *gqlResponse)
}

func (repo *weaviateResponseRepo) prepareNearTextArgumentBuilder(concepts []entities.Concept) *graphql.NearTextArgumentBuilder {

	conceptsStr := make([]string, 0, len(concepts))

	for _, concept := range concepts {
		conceptsStr = append(conceptsStr, string(concept))
	}

	return repo.client.GraphQL().NearTextArgBuilder().
		WithConcepts(conceptsStr)
}

func (repo *weaviateResponseRepo) prepareAskArgBuilder(question string) *graphql.AskArgumentBuilder {

	return repo.client.GraphQL().AskArgBuilder().
		WithQuestion(question)
}

func (repo *weaviateResponseRepo) prepareResponseFromGQLResponse(query entities.Query, gqlResponse graphqlModels.GraphQLResponse) (*entities.Response, error) {

	if err := repo.extractErrorFromGQLResponse(gqlResponse); err != nil {

		return nil, err
	}

	gqlGet, ok := gqlResponse.Data["Get"].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("cannot find root get: %s", gqlResponse.Data)
	}

	gqlDomains, ok := gqlGet[query.DomainId].([]any)

	if !ok {
		return nil, fmt.Errorf("cannot find domains list: %s", gqlGet)
	}

	gqlDomain, _ := gqlDomains[0].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("cannot find domain %s: %s", query.DomainId, gqlDomains)
	}

	addl, _ := gqlDomain["_additional"].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("cannot find additional object: %s", gqlDomain)
	}

	generate, _ := addl["generate"].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("cannot find generated string: %s", addl)
	}

	groupedResult, _ := generate["groupedResult"].(string)

	if !ok {
		return nil, fmt.Errorf("cannot find grouped result: %s", generate)
	}

	source, _ := gqlDomain["source"].(string)

	if !ok {
		return nil, fmt.Errorf("cannot find domains list: %s", gqlDomain)
	}

	return &entities.Response{Query: query, Response: groupedResult, Sources: []string{source}}, nil
}

func (repo *weaviateResponseRepo) extractErrorFromGQLResponse(gqlResponse graphqlModels.GraphQLResponse) error {

	if len(gqlResponse.Errors) < 1 {

		return nil
	}

	errs := make([]string, 0, len(gqlResponse.Errors))

	for _, err := range gqlResponse.Errors {

		errs = append(errs, err.Message)
	}

	return fmt.Errorf("%s", strings.Join(errs, "\n"))
}
