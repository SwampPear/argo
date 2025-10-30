package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"os/exec"
	"bufio"
	"io"
	"encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
	"github.com/SwampPear/argo/pkg/tools"
	"github.com/rs/xid"
)

// Wails app.
type App struct {
	ctx 			 context.Context							 // Wails app context
	projectDir string												 // project directory
	settings	 settings.Settings						 // app settings gathered from scope.yaml
	mu 				 sync.Mutex										 // protects projectDir, settings, and active
  active 		 map[string]context.CancelFunc // maps log event to cancel function
	pw 				 tools.PW											 // Playwright integration
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

		a.mu.Lock()
    if err := json.Unmarshal([]byte(line), &entry); err == nil {
			fmt.Errorf("Failed to process log")
		}

    // broadcast to frontend state
		a.emit(entry)

		a.step++
		a.mu.Unlock()
  }
}

func (a *App) runPipeline(ctx context.Context, id string, cfg settings.Settings) {
	a.mu.Lock()
  start := time.Now()
  log := func(phase, module, action, target, status string, conf float64, d time.Duration, summary string) {
    a.emit(LogEvent{
      Step: a.step, ID: id, Timestamp: time.Now().Format(time.RFC3339), Module: module, Action: action, Target: target, 
			Status: status, Duration: d.String(), Confidence: conf, Summary: summary,
    })
    a.step++
  }
	a.mu.Unlock()

  // 1) Scope/policy check
  // 2) Recon (e.g., run nuclei safe set)
  // 3) Browser crawl (Playwright sidecar or HTTP runner)

  log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
	log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
}

// Creates a new app.
func NewApp() *App {
	return &App{}
}

// Starts an interactive browser.
func (a *App) StartInteractiveBrowser(url string) error {
  cmd := exec.Command("node", "./tools/playwright/pw_runner.mjs", url)
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

// Runs the pipeline.
func (a *App) Run() (string, error) {
  a.mu.Lock()
  cfg := a.settings
  a.mu.Unlock()

  id := xid.New().String()

  ctx, cancel := context.WithCancel(context.Background())

  a.mu.Lock()
  if a.active == nil {
    a.active = make(map[string]context.CancelFunc)
  }
  a.active[id] = cancel
  a.mu.Unlock()

  go a.runPipeline(ctx, id, cfg)
  return id, nil
}

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
