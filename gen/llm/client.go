package llm

type Client interface {
	Call(prompt string) (string, error)
	CallAndCheck(prompt string, check func(string) error) (string, error)
}
