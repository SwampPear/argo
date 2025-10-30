package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"os/exec"
	"os"
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

// Creates a new app.
func NewApp() *App {
	return &App{}
}

// App initialization.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.settings = settings.Default()
}

// Finishes a process determined by log id.
func (a *App) finish(id string) {
	a.mu.Lock()
	if c, ok := a.active[id]; ok {
		c()                    	// cancel the context
		delete(a.active, id) // remove from map
	}
	a.mu.Unlock()
}

func (a *App) StartInteractiveBrowser(url string) error {
  cmd := exec.Command("node", "tools/playwright/pw_runner.mjs", url)
  cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

  // keep child alive while app runs
  return cmd.Start() // don't Wait(); Node stays alive until user closes the browser
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

func (a *App) SelectProjectDirectory() (string, error) {
  dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
    Title: "Select Project Directory",
  })

  if err != nil {
    return "", err
  }

  a.projectDir = dir
	
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

func (a *App) Run() (string, error) {
  a.mu.Lock()
  cfg := a.settings
  a.mu.Unlock()

  runID := xid.New().String()

  ctx, cancel := context.WithCancel(context.Background())
  a.mu.Lock()
  if a.active == nil {
    a.active = make(map[string]context.CancelFunc)
  }
  a.active[runID] = cancel
  a.mu.Unlock()

  go a.runPipeline(ctx, runID, cfg)
  return runID, nil
}

func (a *App) Cancel(runID string) {
  a.mu.Lock(); if c, ok := a.active[runID]; ok { c(); delete(a.active, runID) }; a.mu.Unlock()
}

func (a *App) emit(ev LogEvent) {
  runtime.EventsEmit(a.ctx, "log:event", ev)
}

func (a *App) runPipeline(ctx context.Context, id string, cfg settings.Settings) {
  start := time.Now()
  step := 0
  log := func(phase, module, action, target, status string, conf float64, d time.Duration, summary string) {
    a.emit(LogEvent{
      Step: step, ID: id, Timestamp: time.Now().Format(time.RFC3339), Module: module, Action: action, Target: target, 
			Status: status, Duration: d.String(), Confidence: conf, Summary: summary,
    })
    step++
  }

  // 1) Scope/policy check
	/*
  t0 := time.Now()
  if err := validatePolicy(cfg); err != nil {
    log("Triage","Policy","validate","-", "Error", 0, time.Since(t0), err.Error()); return
  }
  log("Triage","Policy","validate","-", "OK", 1, time.Since(t0), "Policy validated")

  // 2) Recon (e.g., run nuclei safe set)
  t1 := time.Now()
  if err := runNuclei(ctx, cfg, runID, log); err != nil {
    log("Recon","Nuclei","scan","in-scope", "Error", 0, time.Since(t1), err.Error()); return
  }
  log("Recon","Nuclei","scan","in-scope", "OK", 0.5, time.Since(t1), "Low-noise templates completed")

  // 3) Browser crawl (Playwright sidecar or HTTP runner)
  t2 := time.Now()
  if err := crawlPlaywright(ctx, cfg, runID, log); err != nil {
    log("Browser","Playwright","crawl","start_urls", "Error", 0, time.Since(t2), err.Error()); return
  }
  log("Browser","Playwright","crawl","start_urls", "OK", 0.6, time.Since(t2), "HAR/screens captured")

  // 4) LLM triage → hypotheses → targeted tests
  t3 := time.Now()
  findings, err := triageAndTest(ctx, cfg, runID, log)
  if err != nil {
    log("Triage","LLM","cluster+plan","artifacts", "Error", 0, time.Since(t3), err.Error()); return
  }
  log("Test","Prober","validate-hypotheses","endpoints", "OK", 0.8, time.Since(t3), fmt.Sprintf("%d potential issues", len(findings)))
	*/

  log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
	log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
}
