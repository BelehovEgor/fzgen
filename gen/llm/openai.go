package llm

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAiClient struct {
	apiKey string
}

func NewOpenAiClient(apiKey string) *OpenAiClient {
	return &OpenAiClient{
		apiKey: apiKey,
	}
}

var (
	systemAI = `
Respond with "No" or a go-fuzz(https://go.dev/doc/security/fuzz/) test function without any explanation or how-to guide.
Evaluate the given function to check "Is it even worth fuzzing?"
Note: Usually, a good target function does parsing, decoding, deserialization, unmarshaling, etc. of a given input/inputs.
`
)

func (client *OpenAiClient) Call(prompt string) (string, error) {
	if client.apiKey == "" {
		return "", fmt.Errorf("API KEY is required for request.")
	}

	openaiClient := openai.NewClient(
		option.WithAPIKey(client.apiKey),
	)

	ctx := context.Background()

	messages := openai.F([]openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(systemAI + "\n\n" + prompt),
	})

	completion, err := openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       openai.F(openai.ChatModelChatgpt4oLatest),
		Temperature: openai.F(1.000000),
		TopP:        openai.F(1.000000),
	})

	if err != nil {
		return "", err
	}

	if completion.Choices[0].Message.Content != "No" {
		return "", fmt.Errorf("Target was not generated.")
	}

	return completion.Choices[0].Message.Content, nil
}

func (client *OpenAiClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return result, nil
	}

	if checkResult := check(result); checkResult != nil {
		return result, checkResult
	}

	return result, nil
}
