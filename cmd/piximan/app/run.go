package app

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/fekoneko/piximan/internal/resources"
	"github.com/fekoneko/piximan/internal/ui"
)

const applicationId = "com.fekoneko.piximan"

func Run(version string) {
	if len(os.Args) > 1 {
		// TODO: open collection with the path provided
		fmt.Println("providing arguments to the viewer is not yet supported")
	}

	resources.LoadResources()
	runApplication(version)
}

func runApplication(version string) {
	app := adw.NewApplication(applicationId, gio.ApplicationFlagsNone)
	app.SetVersion(version)

	app.ConnectActivate(func() {
		resources.LoadCss()

		window := ui.NewWindow()
		window.Attach(app)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
