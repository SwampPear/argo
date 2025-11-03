package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/state"
	"github.com/SwampPear/argo/pkg/settings"
	"github.com/SwampPear/argo/pkg/runner"
)

// App ties together state, runner, and UI features.
type App struct {
	ctx       context.Context
	stateMgr  *state.Manager
	playwright *runner.Playwright
}

// New creates a new App instance with initialized subsystems.
func New() *App {
	return &App{
		stateMgr:  nil,
		playwright: &runner.Playwright{},
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.stateMgr = state.New(ctx)
}

/* =========================
   State Sync
   ========================= */

func (a *App) GetState() state.AppState {
	return a.stateMgr.GetState()
}

func (a *App) SetState(next state.AppState, baseVersion int64) state.AppState {
	return a.stateMgr.SetState(next, baseVersion)
}

/* =========================
   UI / User Actions
   ========================= */

func (a *App) SelectProjectDirectory() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Directory",
	})
	if err != nil {
		return "", err
	}
	s := a.stateMgr.GetState()
	s.ProjectDir = dir
	a.stateMgr.SetState(s, s.Version)
	return dir, nil
}

func (a *App) LoadYAMLSettings(path string) (settings.Settings, error) {
	cfg, err := settings.LoadYAML(path)
	if err != nil {
		return settings.Settings{}, err
	}
	s := a.stateMgr.GetState()
	s.Settings = cfg
	a.stateMgr.SetState(s, s.Version)
	return cfg, nil
}

/* =========================
   Runner
   ========================= */

func (a *App) StartInteractiveBrowser() error {
	return a.playwright.Start(a.stateMgr)
}

func (a *App) BroadcastState() {
	a.stateMgr.Broadcast()
}
