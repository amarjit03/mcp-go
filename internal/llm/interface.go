package llm

// LLMClient is the interface for LLM providers
type LLMClient interface {
	// Generate sends a prompt and gets a response
	Generate(prompt string, model string) (string, error)

	// IsAvailable checks if the LLM service is available
	IsAvailable() bool

	// GetModels returns available models
	GetModels() ([]string, error)
}
