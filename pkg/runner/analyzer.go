package runner

import (
	"time"

	"github.com/SwampPear/argo/pkg/state"
)

type Analyzer struct{}

// Starts the analyzer.
func (a *Analyzer) Start(m *state.Manager) error {
	m.AppendLog(state.LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Module:     "Analyzer",
		Action:     "start",
		Target:     "analyzer",
		Status:     "ok",
		Duration:   "0ms",
		Confidence: 1.0,
		Summary:    "Analyzer started",
	})

	return nil
}