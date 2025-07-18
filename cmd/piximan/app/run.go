package app

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const applicationId = "com.fekoneko.piximan"
const resourcePrefix = "/com/fekoneko/piximan"

//go:embed piximan.gresource
var resources []byte

func Run(version string) {
	if len(os.Args) > 1 {
		// TODO: open collection with the path provided
		fmt.Println("providing arguments to the viewer is not yet supported")
	}

	registerResources()
	runApplication(version)
}

func registerResources() {
	bytes := glib.NewBytes(resources)
	resource, err := gio.NewResourceFromData(bytes)
	if err != nil {
		panic(err)
	}
	gio.ResourcesRegister(resource)
}

func addCssProvider() {
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromResource(resourcePrefix + "/window.css")
	priority := uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	gtk.StyleContextAddProviderForDisplay(gdk.DisplayGetDefault(), cssProvider, priority)
}

func runApplication(version string) {
	app := adw.NewApplication(applicationId, gio.ApplicationFlagsNone)
	app.SetVersion(version)

	app.ConnectActivate(func() {
		addCssProvider()

		builder := gtk.NewBuilderFromResource(resourcePrefix + "/window.ui")
		window := NewWindow(builder)
		window.Add(app)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
