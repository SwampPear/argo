package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/state"
	"github.com/SwampPear/argo/pkg/settings"
	"github.com/SwampPear/argo/pkg/runner"
)

// Application definition.
type App struct {
	ctx        context.Context
	stateMgr   *state.Manager
	playwright *runner.Playwright
	analyzer   *runner.Analyzer
}

// Initializes a new App instance with new subsystems.
func New() *App {
	return &App{
		stateMgr:   nil,
		playwright: &runner.Playwright{},
		analyzer:	  &runner.Analyzer{},
	}
}

// Called on application startup.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.stateMgr = state.New(ctx)
}

// Gets current state.
func (a *App) GetState() state.RemoteState {
	return a.stateMgr.GetState()
}

// Sets current state.
func (a *App) SetState(next state.RemoteState) state.RemoteState {
	return a.stateMgr.SetState(next)
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
	a.stateMgr.SetState(s)

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
	a.stateMgr.SetState(s)

	return cfg, nil
}

// Begins browser session.
func (a *App) StartInteractiveBrowser() error {
	return a.playwright.Start(a.stateMgr)
}

// Starts the analyzer.
func (a *App) StartAnalyzer() error {
	return a.analyzer.Start(a.stateMgr)
}
