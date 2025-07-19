package ui

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/resources"
)

type Viewer struct {
	*gtk.Box
	window *Window
}

func NewViewer(window *Window) *Viewer {
	builder := resources.NewBuilder("viewer.ui")
	container := builder.GetObject("viewer-container").Cast().(*gtk.Box)

	return &Viewer{container, window}
}

func (v *Viewer) Attach(builder *gtk.Builder) {
	container := builder.GetObject("viewer-root").Cast().(*gtk.ScrolledWindow)
	container.SetChild(v)
}
