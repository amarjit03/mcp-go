package llm

import (
"bytes"
"encoding/json"
"fmt"
"io"
"net/http"
"time"
)

// OllamaClient handles communication with Ollama
type OllamaClient struct {
baseURL string
client  *http.Client
}

// OllamaRequest is the request structure for Ollama
type OllamaRequest struct {
Model  string `json:"model"`
Prompt string `json:"prompt"`
Stream bool   `json:"stream"`
}

// OllamaResponse is the response structure from Ollama
type OllamaResponse struct {
Model    string `json:"model"`
Response string `json:"response"`
Done     bool   `json:"done"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL string) *OllamaClient {
return &OllamaClient{
baseURL: baseURL,
client: &http.Client{
Timeout: 120 * time.Second,
},
}
}

// Generate sends a prompt to Ollama and gets a response
func (oc *OllamaClient) Generate(prompt string, model string) (string, error) {
if model == "" {
model = "tinyllama"
}

request := OllamaRequest{
Model:  model,
Prompt: prompt,
Stream: false,
}

reqBody, err := json.Marshal(request)
if err != nil {
return "", fmt.Errorf("failed to marshal request: %w", err)
}

resp, err := oc.client.Post(
oc.baseURL+"/api/generate",
"application/json",
bytes.NewBuffer(reqBody),
)
if err != nil {
return "", fmt.Errorf("failed to connect to Ollama: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
body, _ := io.ReadAll(resp.Body)
return "", fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
}

var ollamaResp OllamaResponse
if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
return "", fmt.Errorf("failed to decode response: %w", err)
}

return ollamaResp.Response, nil
}

// IsAvailable checks if Ollama is available
func (oc *OllamaClient) IsAvailable() bool {
resp, err := oc.client.Get(oc.baseURL + "/api/tags")
if err != nil {
return false
}
defer resp.Body.Close()
return resp.StatusCode == http.StatusOK
}

// GetModels returns available models
func (oc *OllamaClient) GetModels() ([]string, error) {
resp, err := oc.client.Get(oc.baseURL + "/api/tags")
if err != nil {
return nil, fmt.Errorf("failed to get models: %w", err)
}
defer resp.Body.Close()

var result struct {
Models []struct {
Name string `json:"name"`
} `json:"models"`
}

if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
return nil, fmt.Errorf("failed to decode models: %w", err)
}

var models []string
for _, m := range result.Models {
models = append(models, m.Name)
}

return models, nil
}
