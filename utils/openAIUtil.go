package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	envUtil "service-news-app-backend/config"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func CallLLM(systemPrompt string, message []map[string]interface{}) (map[string]interface{}, error) {

	fmt.Println("systemPrompt: ", systemPrompt)

	messages := []map[string]interface{}{}
	systemRole := map[string]interface{}{
		"role":    "system",
		"content": systemPrompt,
	}

	messages = append(messages, systemRole)

	messages = append(messages, message...)

	// creating input
	llmCompletionCreate := map[string]interface{}{
		"model":       "gpt-3.5-turbo-0125",
		"temperature": 0.3,
		"messages":    messages,
	}

	// make an API call to openAI
	apiKey := envUtil.GetEnvironmentVariable("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewReader(ConvertToJson(llmCompletionCreate)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	responseBody := map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {

		answer := responseBody["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})

		return answer, nil
	} else {
		err := responseBody["error"].(map[string]interface{})["message"].(string)
		return nil, errors.New("can't get response-> " + err)
	}

}

func GenerateVectorEmebeddings(input string) ([]float32, error) {

	apiEndPoint := "https://api.openai.com/v1/embeddings"

	type apiRequest struct {
		Input string `json:"input"`
		Model string `json:"model"`
	}

	data := &apiRequest{
		Input: input,
		Model: "text-embedding-ada-002",
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiEndPoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	apiKey := envUtil.GetEnvironmentVariable("OPENAI_API_KEY")

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(result["error"].(map[string]interface{})["message"].(string))
	} else {

		resultEmbeddings := result["data"].([]interface{})[0].(map[string]interface{})["embedding"].([]interface{})

		var embeddings []float32
		for _, v := range resultEmbeddings {
			embeddings = append(embeddings, float32(v.(float64)))
		}

		return embeddings, nil
	}

}

type Entities struct {
	Organizations string   `json:"organizations"`
	Locations     []string `json:"locations"`
	Individuals   []string `json:"individuals"`
}

// TaggingData represents the structure of the response payload
type MetaData struct {
	SentimentScore string   `json:"sentimentScore"`
	Entities       Entities `json:"entities"`
	Categories     []string `json:"categories"`
}

// getCategories calls the OpenAI API to get categories based on the prompt
func GetResponseFromChatGPT(ctx context.Context, content string) (*MetaData, error) {

	// Create the prompt for OpenAI
	prompt := "Extract entities from the following news article and categorize them into organizations, locations, and individuals:\n\n" + content

	req := openai.CompletionRequest{
		Model:       openai.GPT4o,
		Prompt:      prompt,
		MaxTokens:   50,
		Temperature: 0.3,
	}

	// Initialize OpenAI client
	var openAIClient *openai.Client

	resp, err := openAIClient.CreateCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no choices returned from OpenAI")
	}

	responseFromOpenAI := strings.TrimSpace(resp.Choices[0].Text)

	// {
	// 	"sentimentScore": "Neutral",
	// 	"categories": ["National Security", "Conflict"],
	// 	"entities": {
	// 		"organizations": ["Example Publisher"],
	// 		"locations": [],
	// 		"individuals": []
	// 	}
	// }

	var resultData *MetaData

	// Unmarshal the JSON into the MetaData struct
	err = json.Unmarshal([]byte(responseFromOpenAI), &resultData)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return resultData, nil
}

func GenerateSummary(ctx context.Context, content string) (string, error) {

	// Create the prompt for OpenAI
	prompt := "Extract entities from the following news article and categorize them into organizations, locations, and individuals:\n\n" + content
	req := openai.CompletionRequest{
		Model:       openai.GPT4o,
		Prompt:      prompt,
		MaxTokens:   50,
		Temperature: 0.3,
	}

	// Initialize OpenAI client
	var openAIClient *openai.Client

	resp, err := openAIClient.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned from OpenAI")
	}

	responseFromOpenAI := strings.TrimSpace(resp.Choices[0].Text)

	return responseFromOpenAI, nil
}

// parseCategories attempts to parse the categories as JSON
func ParseCategories(categoriesText string) ([]string, error) {
	var categories []string
	err := json.Unmarshal([]byte(categoriesText), &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func GenerateEmbeddings(ctx context.Context, prompt string) (string, error) {
	req := openai.CompletionRequest{
		Model:       openai.GPT4o,
		Prompt:      prompt,
		MaxTokens:   50,
		Temperature: 0.3,
	}

	// Initialize OpenAI client
	var openAIClient *openai.Client

	resp, err := openAIClient.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices returned from OpenAI")
	}

	return strings.TrimSpace(resp.Choices[0].Text), nil
}
