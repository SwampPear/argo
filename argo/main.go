package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)


var assets embed.FS // go:embed all:frontend/dist

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Argo",
		DisableResize: false,
		Width:  1024,
		Height: 768,
		MinWidth:    400,
    MinHeight:   300,
    MaxWidth:    0,
    MaxHeight:   0,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
