package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/SwampPear/argo/pkg/settings"
)

type App struct {
	ctx 			 context.Context
	projectDir string
	settings	 settings.Settings
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