package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/state"
	"github.com/SwampPear/argo/pkg/settings"
	"github.com/SwampPear/argo/pkg/runner"
)

// App.
type App struct {
	ctx       context.Context
	stateMgr  *state.Manager
	playwright *runner.Playwright
}

// Creates a new App instance with initialized subsystems.
func New() *App {
	return &App{
		stateMgr:  nil,
		playwright: &runner.Playwright{},
	}
}

// Called on application startup.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.stateMgr = state.New(ctx)
}

// Gets current state.
func (a *App) GetState() state.AppState {
	return a.stateMgr.GetState()
}

// Sets current state.
func (a *App) SetState(next state.AppState, baseVersion int64) state.AppState {
	return a.stateMgr.SetState(next, baseVersion)
}

// Hydrates frontend with current state.
func (a *App) BroadcastState() {
	a.stateMgr.Broadcast()
}

// Selects project directory and updates state.
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

// Loads settings from <projectDir/scope.yaml and updates state.
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

// Begins browser session.
func (a *App) StartInteractiveBrowser() error {
	return a.playwright.Start(a.stateMgr)
}

