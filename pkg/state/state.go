package state

import (
	"context"
	"sync"
	"strings"
	"sort"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
)

// Log entry for tool calls and info.
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

// Shared remote state.
type RemoteState struct {
	ProjectDir  string            `json:"project_dir"`
	Settings    settings.Settings `json:"settings"`
	Logs        []LogEntry        `json:"logs"`
	ScopeFilter bool							`json:"scope_filter"`
}

// State manager.
type Manager struct {
	mu      sync.RWMutex
	ctx     context.Context
	state   RemoteState
}

// Initializes a new state manager.
func New(ctx context.Context) *Manager {
	return &Manager{
		ctx: ctx,
		state: RemoteState{
			ProjectDir:	 "",
			Settings: 	 settings.Default(),
			Logs:     	 make([]LogEntry, 0, 256),
			ScopeFilter: false,
		},
	}
}

// Safely returns current state snapshot.
func (m *Manager) GetState() RemoteState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.state
}

// Safely applies a new state.
func (m *Manager) SetState(next RemoteState) RemoteState {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state = next

	runtime.EventsEmit(m.ctx, "state:update", next)

	return next
}

// Hydrates the frontend with current state.
func (m *Manager) Broadcast() {
	m.mu.RLock()
	s := m.state
	m.mu.RUnlock()

	runtime.EventsEmit(m.ctx, "state:update", s)
}

// Adds a log entry and emits an update event.
func (m *Manager) AppendLog(le LogEntry) {
	m.mu.Lock()
	m.state.Logs = append(m.state.Logs, le)
	m.mu.Unlock()

	runtime.EventsEmit(m.ctx, "log:event", le)
}

// Gets logs for the analyzer.
func (m *Manager) Logs() []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// filter
	scopeFilter := m.state.ScopeFilter
	
	// array of the scopes from the state settings
	assets := m.state.Settings.Assets.InScope
	n := len(assets)
	scopes := make([]string, 0, n)
	for i := 0; i < n; i++ {
		scopes = append(scopes, assets[i].Hostname)
	}

	// filter logs
	out := make([]LogEntry, 0, len(m.state.Logs))
	for _, e := range m.state.Logs {
		// by type
		if strings.EqualFold(strings.TrimSpace(e.Module), "Analyzer") {
			continue
		}

		// by scope
		if scopeFilter {
			t := strings.TrimSpace(e.Target)
			match := false
			for _, host := range scopes {
					if strings.Contains(t, host) {
							match = true
							break
					}
			}
			if !match {
					continue
			}
		}

		out = append(out, e)
	}

	// order logs by time
	sort.SliceStable(out, func(i, j int) bool {
		ti, tj := out[i].Timestamp, out[j].Timestamp
		if ti != tj {
			return ti < tj
		}

		return false
	})

	return out
}