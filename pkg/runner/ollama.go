package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OllamaClient struct {
	BaseURL 		string
	Model   		string
	Timeout     time.Duration
	Temperature float64
}

func (c *OllamaClient) Complete(ctx context.Context, system, prompt string) (string, error) {
	type msg struct{ Role, Content string }
	body := map[string]any{
		"model": c.Model,
		"messages": []msg{
			{Role: "system", Content: system},
			{Role: "user", Content: prompt},
		},
		"stream":      false,
		"temperature": c.Temperature,
	}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/chat", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: ifZero(c.Timeout, 60*time.Second)}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Message struct{ Content string } `json:"message"`
		Error   string                   `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Error != "" {
		return "", fmt.Errorf(out.Error)
	}
	
	return out.Message.Content, nil
}

func ifZero[T comparable](v, d T) T { var z T; if v == z { return d }; return v }
