package ui

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/fekoneko/piximan/internal/resources"
)

type Window struct {
	*adw.ApplicationWindow
}

func NewWindow() *Window {
	builder := resources.NewBuilder("window.ui")
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
