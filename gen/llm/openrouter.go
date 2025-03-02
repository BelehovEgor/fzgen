package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type OpenRouterClient struct {
}

func (client *OpenRouterClient) Call(prompt string) (string, error) {
	baseURL := "https://openrouter.ai/api/v1"
	apiKey := "your key"

	// Define the request payload
	payload := map[string]interface{}{
		"model": "qwen/qwen2.5-vl-72b-instruct:free",
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
				},
			},
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %s", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %s", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %s", err.Error())
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %s", err.Error())
	}

	time.Sleep(1 * time.Second)

	// Extract and print the completion content
	choices := result["choices"].([]interface{})
	if len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		message := firstChoice["message"].(map[string]interface{})
		content := message["content"].(string)

		content = strings.ReplaceAll(content, "```go", "")
		content = strings.ReplaceAll(content, "```", "")

		return content, nil
	}

	return "", fmt.Errorf("no choices found in the response.")
}

func (client *OpenRouterClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return "", nil
	}

	if err = check(result); err != nil {
		return "", nil
	}

	return result, nil
}
