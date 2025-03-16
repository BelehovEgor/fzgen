package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type GroqDeepSeekClient struct {
}

func (client *GroqDeepSeekClient) Call(prompt string) (string, error) {
	apiKey := os.Getenv("groq_api_key")
	if apiKey == "" {
		return "", fmt.Errorf("no api key")
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	requestBody := ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "you are a professional golang developer who completes fuzzing tests",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:  "deepseek-r1-distill-llama-70b",
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Error marshaling request body:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	http_client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := http_client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request:", err)
	}
	defer resp.Body.Close()

	timeout, err := strconv.Atoi(os.Getenv("timeout"))
	if err == nil || timeout > 0 {
		time.Sleep(time.Duration(timeout) * time.Second)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body:", err)
	}

	var chatCompletionResponse ChatCompletionResponse
	err = json.Unmarshal(body, &chatCompletionResponse)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling response body:", err)
	}

	if len(chatCompletionResponse.Choices) > 0 {
		content := chatCompletionResponse.Choices[0].Message.Content

		content = strings.ReplaceAll(content, "```go", "")
		content = strings.ReplaceAll(content, "```", "")

		content = strings.ReplaceAll(content, "<think>", "/*<think>")
		content = strings.ReplaceAll(content, "</think>", "</think>*/")

		return content, nil
	} else {
		fmt.Println("No choices returned in the response")
	}

	return "", fmt.Errorf("no choices found in the response.")
}

func (client *GroqDeepSeekClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return "", nil
	}

	if err = check(result); err != nil {
		return "", nil
	}

	return result, nil
}
