package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/SwampPear/argo/pkg/settings"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	Version    int64             `json:"version"` // <-- embed version
}

type App struct {
	ctx    context.Context
	mu     sync.RWMutex
	state  AppState
	// keep an internal counter for mutations; we mirror it into state.Version before returning/emitting
	version int64

	active map[string]context.CancelFunc
	step   int
}

func NewApp() *App {
	return &App{
		state: AppState{
			ProjectDir: "",
			Settings:   settings.Default(),
			Logs:       make([]LogEntry, 0, 256),
			Version:    0,
		},
		version: 0,
		active:  make(map[string]context.CancelFunc),
		step:    0,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

/* =========================
   State Sync API
   ========================= */

// Always return a single JSON object
func (a *App) GetState() AppState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	s := a.state
	s.Version = a.version
	return s
}

// Return authoritative state as a single JSON object; client infers acceptance by comparing versions
func (a *App) SetState(next AppState, baseVersion int64) AppState {
	a.mu.Lock()
	defer a.mu.Unlock()

	if baseVersion != a.version {
		// reject: return canonical current state
		s := a.state
		s.Version = a.version
		return s
	}

	a.state = next
	a.version++
	s := a.state
	s.Version = a.version

	// broadcast a single payload
	runtime.EventsEmit(a.ctx, "state:update", s)
	return s
}

// --- internal helper; caller holds a.mu
func (a *App) appendLogUnsafe(le LogEntry) {
	if le.Step == 0 {
		a.step++
		le.Step = a.step
	} else if le.Step > a.step {
		a.step = le.Step
	}
	a.state.Logs = append(a.state.Logs, le)
	a.version++
}

func (a *App) AppendLog(le LogEntry) {
	a.mu.Lock()
	a.appendLogUnsafe(le)
	s := a.state
	s.Version = a.version
	a.mu.Unlock()

	// single payload
	runtime.EventsEmit(a.ctx, "state:update", s)
	// (optional legacy)
	runtime.EventsEmit(a.ctx, "log:event", le)
}

/* =========================
   Features that mutate state
   ========================= */

func (a *App) SelectProjectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Directory",
	})
	if err != nil {
		return "", err
	}

	a.mu.Lock()
	a.state.ProjectDir = dir
	a.version++
	s := a.state
	s.Version = a.version
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "state:update", s)
	return dir, nil
}

func (a *App) LoadYAMLSettings(path string) (settings.Settings, error) {
	cfg, err := settings.LoadYAML(path)
	if err != nil {
		return settings.Settings{}, fmt.Errorf("load settings: %w", err)
	}

	a.mu.Lock()
	a.state.Settings = cfg
	a.version++
	s := a.state
	s.Version = a.version
	a.mu.Unlock()

	runtime.EventsEmit(a.ctx, "state:update", s)
	return cfg, nil
}

/* =========================
   Playwright Runner + Log Streaming
   ========================= */

func (a *App) StartInteractiveBrowser() error {
	cmd := exec.Command("node", "./tools/playwright/pw_runner.mjs")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go a.streamLogs(stdout)
	go a.streamLogs(stderr)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (a *App) StopInteractiveBrowser(id string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if c, ok := a.active[id]; ok {
		c()
		delete(a.active, id)
		return true
	}
	return false
}

func (a *App) streamLogs(r io.Reader) {
	sc := bufio.NewScanner(r)
	buf := make([]byte, 0, 1024*64)
	sc.Buffer(buf, 1024*1024)

	for sc.Scan() {
		line := sc.Text()
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry = LogEntry{
				Summary:   line,
				Status:    "OK",
				Module:    "Playwright",
				Action:    "log",
			}
		}
		a.AppendLog(entry)
	}
}

/* =========================
   Utility
   ========================= */

func (a *App) BroadcastState() {
	a.mu.RLock()
	s := a.state
	v := a.version
	a.mu.RUnlock()
	s.Version = v
	runtime.EventsEmit(a.ctx, "state:update", s)
}
