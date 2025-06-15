package app

import (
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const APPLICATION_ID = "com.fekoneko.piximan"
const WINDOW_TITLE = "piximan"

func Run() {
	app := adw.NewApplication(APPLICATION_ID, 0)
	app.ConnectActivate(func() { activate(&app.Application) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	header := adw.NewHeaderBar()
	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(header)

	window := adw.NewApplicationWindow(app)
	window.SetDefaultSize(600, 300)
	window.SetTitle(WINDOW_TITLE)
	window.SetContent(box)
	window.SetVisible(true)
}
