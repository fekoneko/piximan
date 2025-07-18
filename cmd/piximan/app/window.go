package app

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Window struct {
	*adw.ApplicationWindow
}

func NewWindow() *Window {
	builder := gtk.NewBuilderFromResource(resourcePrefix + "/window.ui")
	window := builder.GetObject("window").Cast().(*adw.ApplicationWindow)

	explorer := NewExplorer()
	explorer.Attach(builder)

	viewer := NewViewer()
	viewer.Attach(builder)

	return &Window{window}
}

func (w *Window) Attach(app *adw.Application) {
	app.AddWindow(&w.Window)
}
