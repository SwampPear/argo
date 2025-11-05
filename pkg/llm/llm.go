package llm

import (
	"context"
)

// Abstract client interface.
type LLMClient interface {
	Complete(ctx context.Context, system, prompt string) (string, error)
}