package app

import (
	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Window struct {
	*adw.ApplicationWindow
	worksList *WorksList
}

func NewWindow(builder *gtk.Builder) *Window {
	window := builder.GetObject("window").Cast().(*adw.ApplicationWindow)
	worksList := NewWorksList(builder)

	return &Window{window, worksList}
}

func (w *Window) Add(app *adw.Application) {
	app.AddWindow(&w.Window)
}
