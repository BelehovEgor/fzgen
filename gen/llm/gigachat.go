package llm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	exists_bearer_token string
)

type GigachatClient struct {
}

func (client *GigachatClient) Call(prompt string) (string, error) {
	token, err := createToken()
	if err != nil {
		return "", err
	}

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	method := "POST"

	payload := map[string]interface{}{
		"model": "GigaChat-Max",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "Ты профессиональный разработчик на языке golang, который пишет фаззинг тесты. Расширь фаззинг цели.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"stream":          false,
		"update_interval": 0,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %s", err)
	}

	http_client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := http_client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil || result == nil {
		return "", fmt.Errorf("error unmarshaling response: %s", err.Error())
	}

	time.Sleep(1 * time.Second)

	// Extract and print the completion content
	choises_interface, ok := result["choices"]
	if !ok {
		return "", fmt.Errorf("error getting result")
	}

	choices := choises_interface.([]interface{})

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

func (client *GigachatClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return "", nil
	}

	if err = check(result); err != nil {
		return "", nil
	}

	return result, nil
}

func createToken() (string, error) {
	if exists_bearer_token != "" {
		return exists_bearer_token, nil
	}

	// Define the URL
	auth_url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"

	// Define the payload
	payload := url.Values{}
	payload.Set("scope", "GIGACHAT_API_PERS")

	// Define the headers
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Accept":        "application/json",
		"RqUID":         "22eaee32-64a2-4199-b8c2-041d4c1c809d",
		"Authorization": "Basic " + os.Getenv("gigachat_api_key"),
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", auth_url, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil || result == nil {
		return "", fmt.Errorf("error unmarshaling response: %s", err.Error())
	}

	bearer_token, ok := result["access_token"]
	if !ok {
		return "", fmt.Errorf("failed")
	}

	exists_bearer_token = bearer_token.(string)

	return exists_bearer_token, nil
}
