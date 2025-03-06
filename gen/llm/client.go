package llm

import "fmt"

type Client interface {
	Call(prompt string) (string, error)
	CallAndCheck(prompt string, check func(string) error) (string, error)
}

type MockClient struct {
}

func (client *MockClient) Call(prompt string) (string, error) {
	return "", fmt.Errorf("it is mock client")
}

func (client *MockClient) CallAndCheck(prompt string, check func(string) error) (string, error) {
	result, err := client.Call(prompt)
	if err != nil {
		return "", nil
	}

	if err = check(result); err != nil {
		return "", nil
	}

	return result, nil
}

func GetClient(clientName string) Client {
	switch clientName {
	case "openrouter":
		return &OpenRouterClient{}
	case "gigachat":
		return &GigachatClient{}
	case "mock":
		return &MockClient{}
	default:
		panic("unimplemented client")
	}
}
