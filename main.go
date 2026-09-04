package main

import (
	"embed"
	"log"
	"time"
	"os/exec"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var embeddedFiles embed.FS

func init() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[string]("tunnel-log")
}

func main() {
	// Create the service first.
	tunnelService := &TunnelService{
		procs: make(map[int]*exec.Cmd),
	}

	// Create the Wails application.
	app := application.New(application.Options{
		Name:        "TunnelPilot v2.1",
		Description: "TunnelPilot",

		Services: []application.Service{
			application.NewService(tunnelService),
		},

		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(embeddedFiles),
		},

		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Give TunnelService access to the Wails event manager.
	tunnelService.event = app.Event

	// Create window.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:    "TunnelPilot",
		Width:    1024,
		Height:   668,
		MinWidth: 1024,
		MinHeight:668,

		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},

		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	// Example time event.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)

			app.Event.Emit("time", now)

			time.Sleep(time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}