package state

import (
	"context"
	"sync"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
)

type LogEntry struct {
	Step         int     `json:"step"`
	ID           string  `json:"id"`
	Timestamp    string  `json:"timestamp"`
	Module       string  `json:"module"`
	Action       string  `json:"action"`
	Target       string  `json:"target"`
	Status       string  `json:"status"`
	Duration     string  `json:"duration"`
	Confidence   float64 `json:"confidence"`
	Summary      string  `json:"summary"`
	ParentStepID int     `json:"parent_step_id"`
}

type AppState struct {
	ProjectDir string            `json:"projectDir"`
	Settings   settings.Settings `json:"settings"`
	Logs       []LogEntry        `json:"logs"`
	Version    int64             `json:"version"`
}

type Manager struct {
	mu      sync.RWMutex
	ctx     context.Context
	state   AppState
	version int64
	step    int
}

func New(ctx context.Context) *Manager {
	return &Manager{
		ctx: ctx,
		state: AppState{
			Settings: settings.Default(),
			Logs:     make([]LogEntry, 0, 256),
		},
	}
}

// GetState safely returns the current immutable state snapshot.
func (m *Manager) GetState() AppState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s := m.state
	s.Version = m.version
	return s
}

// SetState applies a new version if the base version matches.
func (m *Manager) SetState(next AppState, baseVersion int64) AppState {
	m.mu.Lock()
	defer m.mu.Unlock()

	if baseVersion != m.version {
		s := m.state
		s.Version = m.version
		return s
	}
	m.state = next
	m.version++
	s := m.state
	s.Version = m.version
	runtime.EventsEmit(m.ctx, "state:update", s)
	return s
}

// Hydrates the frontend with current state.
func (m *Manager) Broadcast() {
	m.mu.RLock()
	s := m.state
	v := m.version
	m.mu.RUnlock()

	s.Version = v
	runtime.EventsEmit(m.ctx, "state:update", s)
}

// AppendLog adds a log entry and emits an update event.
func (m *Manager) AppendLog(le LogEntry) {
	m.mu.Lock()
	if le.Step == 0 {
		m.step++
		le.Step = m.step
	}
	m.state.Logs = append(m.state.Logs, le)
	m.version++
	s := m.state
	s.Version = m.version
	m.mu.Unlock()

	runtime.EventsEmit(m.ctx, "state:update", s)
	runtime.EventsEmit(m.ctx, "log:event", le)
}

// Gets logs for the analyzer.
func (m *Manager) Logs() []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]LogEntry, 0, len(m.state.Logs))
	for _, e := range m.state.Logs {
		if strings.EqualFold(strings.TrimSpace(e.Module), "Analyzer") {
			continue
		}
		
		out = append(out, e)
	}
	return out
}