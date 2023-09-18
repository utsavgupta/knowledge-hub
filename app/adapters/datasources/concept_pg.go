package datasources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/services"
)

const openaiChatURL = "https://api.openai.com/v1/chat/completions"

type requestBody struct {
	Model       string              `json:"model"`
	Messages    []map[string]string `json:"messages"`
	Temperature int                 `json:"temperature"`
}

type responseBody struct {
	Choices []map[string]any `json:"choices"`
}

type conceptOpenAI struct {
	httpClient      *http.Client
	openaiAccessKey string
}

func NewConceptOpenAI(httpClient *http.Client, openaiAccessKey string) services.ConceptService {

	return &conceptOpenAI{httpClient, openaiAccessKey}
}

func (service *conceptOpenAI) Get(ctx context.Context, question string) ([]entities.Concept, error) {

	requestBody := service.prepareRequestBody(question)
	bRequestBody, err := service.convertRequestBodyToBytes(requestBody)

	if err != nil {
		return nil, err
	}

	responseBody, err := service.makeRequest(bRequestBody)

	if err != nil || responseBody == nil {
		return nil, err
	}

	return service.extractConcepts(*responseBody)
}

func (service *conceptOpenAI) makeRequest(body []byte) (*responseBody, error) {

	requestBodyReader := bytes.NewReader(body)
	request, err := http.NewRequest(http.MethodPost, openaiChatURL, requestBodyReader)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", service.openaiAccessKey))

	if err != nil {
		return nil, fmt.Errorf("could not create request object for Open AI: %w", err)
	}

	httpResponse, err := service.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("could not complete request to Open AI: %w", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Open AI sent back status code %d", httpResponse.StatusCode)
	}

	defer httpResponse.Body.Close()

	var response responseBody
	err = json.NewDecoder(httpResponse.Body).Decode(&response)

	if err != nil {
		return nil, fmt.Errorf("could not parse response received from Open AI: %w", err)
	}

	return &response, nil
}

func (service *conceptOpenAI) extractConcepts(body responseBody) ([]entities.Concept, error) {

	if body.Choices == nil || len(body.Choices) < 1 {
		return nil, nil
	}

	message, ok := body.Choices[0]["message"]

	if !ok {
		return nil, fmt.Errorf("`message` property was not found in the response body received from Open AI")
	}

	messageMap, ok := message.(map[string]any)

	if !ok {
		return nil, fmt.Errorf("expected `message` property to be an object")
	}

	content, ok := messageMap["content"]

	if !ok {
		return nil, fmt.Errorf("`message.content` property was not found in the response body received from Open AI")
	}

	contentStr, ok := content.(string)

	if !ok {
		return nil, fmt.Errorf("expected `message.content` to be a string")
	}

	var concepts []entities.Concept

	if err := json.Unmarshal([]byte(contentStr), &concepts); err != nil {
		return nil, fmt.Errorf("could not unmarshal message content %v: %w", content, err)
	}

	return concepts, nil
}

func (service *conceptOpenAI) convertRequestBodyToBytes(body requestBody) ([]byte, error) {

	b, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("could not marshall request body %v: %w", body, err)
	}

	return b, err
}

func (service *conceptOpenAI) prepareRequestBody(question string) requestBody {

	messages := make([]map[string]string, 0, 2)
	messages = append(messages, service.prepareSystemMessage())
	messages = append(messages, service.prepareUserMessage(question))

	return requestBody{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0,
	}
}

func (service *conceptOpenAI) prepareSystemMessage() map[string]string {

	message := make(map[string]string)
	message["role"] = "system"
	message["content"] = "Extract the concepts from the user's question and return the results in the form of an array. Escape the strings in the response."

	return message
}

func (service *conceptOpenAI) prepareUserMessage(question string) map[string]string {

	message := make(map[string]string)
	message["role"] = "user"
	message["content"] = question

	return message
}
