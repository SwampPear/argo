package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
	"github.com/rs/xid"
)

type App struct {
	ctx 			 context.Context
	projectDir string
	settings	 settings.Settings
	mu sync.Mutex
  active map[string]context.CancelFunc
}

type LogEvent struct {
  Timestamp string `json:"timestamp"`
  RunID     string `json:"run_id"`
  StepID    string `json:"step_id"`
  Module    string `json:"module"`
  Action    string `json:"action"`
  Target    string `json:"target"`
  Status    string `json:"status"`
  Duration  string `json:"duration"`
  Confidence float64 `json:"confidence"`
  Summary   string `json:"summary"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.settings = settings.Default()
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

  runID := time.Now().UTC().Format("20060102T150405Z") + "-" + xid.New().String()

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
  // Wails event to TS: listen via EventsOn('log:event', ...)
  runtime.EventsEmit(a.ctx, "log:event", ev)
}

func (a *App) runPipeline(ctx context.Context, runID string, cfg settings.Settings) {
  start := time.Now()
  step := 0
  log := func(phase, module, action, target, status string, conf float64, d time.Duration, summary string) {
    a.emit(LogEvent{
      Timestamp: time.Now().Format(time.RFC3339),
      RunID: runID, StepID: fmt.Sprintf("step-%03d", step), Module: module,
      Action: action, Target: target, Status: status, Duration: d.String(), Confidence: conf, Summary: summary,
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

  // 5) Report draft
  t4 := time.Now()
  if err := draftReports(ctx, cfg, runID, findings, log); err != nil {
    log("Report","LLM","draft","issues", "Error", 0, time.Since(t4), err.Error()); return
  }
  log("Report","LLM","draft","issues", "OK", 0.9, time.Since(t4), "Drafts generated")
	*/

  log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
}
