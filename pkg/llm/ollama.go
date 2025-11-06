package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"github.com/SwampPear/argo/pkg/state"
)

// Ollama LLM client.
type OllamaClient struct {
	BaseURL     string  		  // local api root url ("http://localhost:11434")
	Model       string				// model name
	Timeout     time.Duration // e120 * time.Second
	Temperature float64       // e.g. 0.4
}

// Initializes Ollama llm client.
func (c *OllamaClient) Init(m *state.Manager) error {
	cfg := m.GetState().Settings.LLM

	c = &OllamaClient{
		BaseURL:     cfg.BaseURL,
		Model:       cfg.Model,
		Temperature: cfg.Temperature,
		Timeout:     time.Duration(cfg.Timeout),
	}

	return nil
}

func (c *OllamaClient) Complete(ctx context.Context, system, prompt string) (string, error) {
	// request body
	type msg struct{ Role, Content string }
	payload := map[string]any{
		"model":   c.Model,
		"messages": []msg{
			{Role: "system", Content: system},
			{Role: "user", Content: prompt},
		},
		"stream": false,
		"options": map[string]any{ // Ollama expects sampling under options
			"temperature": c.Temperature,
		},
	}

	// json parse request body
	b, err := json.Marshal(payload)
	if err != nil { return "", err }

	// request url
	chatURL, err := url.JoinPath(c.BaseURL, "/api/chat")
	if err != nil { return "", err }

	// make request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatURL, bytes.NewReader(b))
	if err != nil { return "", err }
	req.Header.Set("Content-Type", "application/json")

	// send request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil { return "", err }
	defer resp.Body.Close()

	// error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama %s: %s", resp.Status, string(body))
	}

	// output
	var out struct {
		Message struct{ Content string } `json:"message"`
		Error   string                   `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil { return "", err }
	if out.Error != "" { return "", fmt.Errorf(out.Error) }

	return out.Message.Content, nil
}