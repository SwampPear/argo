package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/SwampPear/argo/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	a := app.New() // use internal/app.New()

	err := wails.Run(&options.App{
		Title:            "Argo",
		DisableResize:    false,
		WindowStartState: options.Normal,
		Mac: &mac.Options{
			TitleBar: mac.TitleBarDefault(),
		},
		Width:     1024,
		Height:    768,
		MinWidth:  400,
		MinHeight: 300,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: a.Startup, // note: exported Startup method
		Bind: []interface{}{
			a,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
