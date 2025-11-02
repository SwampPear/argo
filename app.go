package main

import (
	"context"
	"fmt"
	"sync"
	"os/exec"
	"bufio"
	"io"
	"encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
	//"github.com/rs/xid"
)

// Wails app.
type App struct {
	ctx 			 context.Context							 // Wails app context
	projectDir string												 // project directory
	settings	 settings.Settings						 // app settings gathered from scope.yaml
	mu 				 sync.Mutex										 // protects projectDir, settings, and active
  active 		 map[string]context.CancelFunc // maps log event to cancel function
	step			 int													 // current log step
}

// Log events.
type LogEvent struct {
	Step			 int		 `json:"step"`
	ID    		 string	 `json:"id"`
  Timestamp  string  `json:"timestamp"`
  Module     string  `json:"module"`
  Action     string  `json:"action"`
  Target     string  `json:"target"`
  Status     string  `json:"status"`
  Duration   string  `json:"duration"`
  Confidence float64 `json:"confidence"`
  Summary    string  `json:"summary"`
}

// App initialization.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.settings = settings.Default()
}

// Emits a log event.
func (a *App) emit(ev LogEvent) {
  runtime.EventsEmit(a.ctx, "log:event", ev)
}

// Streams logs to frontend and terminal.
func (a *App) streamLogs(r io.Reader) {
  sc := bufio.NewScanner(r)
  buf := make([]byte, 0, 1024*64)
  sc.Buffer(buf, 1024*1024)

  for sc.Scan() {
    line := sc.Text()

		fmt.Println(line)

    var entry LogEvent
    if err := json.Unmarshal([]byte(line), &entry); err == nil {
			fmt.Errorf("Failed to process log")
		}

		a.mu.Lock()
		a.emit(entry)
		a.step++
		a.mu.Unlock()
  }
}

// Creates a new app.
func NewApp() *App {
	return &App{}
}

// Starts an interactive browser.
func (a *App) StartInteractiveBrowser() error {
	// playwright node runner
  cmd := exec.Command("node", "./tools/playwright/pw_runner.mjs")

	// stream
  stdout, _ := cmd.StdoutPipe()
  stderr, _ := cmd.StderrPipe()
  go a.streamLogs(stdout)
  go a.streamLogs(stderr)

  if err := cmd.Start(); err != nil { return err }
  
  return nil
}

// Stops an interactive browser session.
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

// Selects project directory from file folder and sets in app state.
func (a *App) SelectProjectDirectory() (string, error) {
  dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
    Title: "Select Project Directory",
  })

  if err != nil {
    return "", err
  }

	a.mu.Lock()
  a.projectDir = dir
	a.mu.Unlock()
	
  return dir, nil
}

// Loads settings from YAML file.
func (a *App) LoadYAMLSettings(path string) (settings.Settings, error) {
  cfg, err := settings.LoadYAML(path)
  if err != nil {
    return settings.Settings{}, fmt.Errorf("load settings: %w", err)
  }

  a.mu.Lock()
  a.settings = cfg
  a.mu.Unlock()

  return cfg, nil
}