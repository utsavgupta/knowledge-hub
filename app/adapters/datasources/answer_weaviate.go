package datasources

import (
	"context"
	"fmt"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/repos"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	graphqlModels "github.com/weaviate/weaviate/entities/models"
)

type weaviateAnswerRepo struct {
	client *weaviate.Client
}

func NewWeaviateAnswerRepo(scheme string, host string, openaiAccessKey string) (repos.AnswerRepo, error) {

	cfg := weaviate.Config{
		Host:    host,
		Scheme:  scheme,
		Headers: map[string]string{"X-OpenAI-Api-Key": openaiAccessKey},
	}

	client, err := weaviate.NewClient(cfg)

	if err != nil {
		return nil, fmt.Errorf("could not create Weaviate client with config %v: %w", cfg, err)
	}

	return &weaviateAnswerRepo{client}, nil
}

func (repo *weaviateAnswerRepo) Get(ctx context.Context, query entities.Query) (*entities.Answer, error) {

	generativeSearchBuilder := graphql.NewGenerativeSearch().GroupedResult(query.Question)
	nearTextArgumentBuilder := repo.prepareNearTextArgumentBuilder(query.Concepts)
	// fields := []graphql.Field{
	// 	{Name: "title"},
	// 	{Name: "_additional", Fields: []graphql.Field{
	// 		{Name: "answer", Fields: []graphql.Field{
	// 			{Name: "hasAnswer"},
	// 			{Name: "property"},
	// 			{Name: "result"},
	// 			{Name: "startPosition"},
	// 			{Name: "endPosition"},
	// 		}},
	// 	}},
	// }

	fmt.Printf("%v\n", query)

	q := repo.client.GraphQL().
		Get().
		WithClassName(query.DomainId).
		// WithFields(graphql.Field{Name: "text"}, graphql.Field{Name: "source"}).
		WithGenerativeSearch(generativeSearchBuilder).
		WithNearText(nearTextArgumentBuilder).
		WithLimit(3)

	response, err := q.Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve answer from Weaviate for question `%s`: %w", query.Question, err)
	}

	fmt.Printf("%v\n", response.Data)

	for _, errs := range response.Errors {
		fmt.Printf("%v\n", errs)
	}

	return nil, nil
}

func (repo *weaviateAnswerRepo) prepareNearTextArgumentBuilder(concepts []entities.Concept) *graphql.NearTextArgumentBuilder {

	conceptsStr := make([]string, 0, len(concepts))

	for _, concept := range concepts {
		conceptsStr = append(conceptsStr, string(concept))
	}

	return repo.client.GraphQL().NearTextArgBuilder().
		WithConcepts(conceptsStr)
}

func (repo *weaviateAnswerRepo) prepareAskArgBuilder(question string) *graphql.AskArgumentBuilder {

	return repo.client.GraphQL().AskArgBuilder().
		WithQuestion(question)
}

func (repo *weaviateAnswerRepo) prepareAnswerFromResponse(query entities.Query, response graphqlModels.GraphQLResponse) *entities.Answer {
	return nil
}
