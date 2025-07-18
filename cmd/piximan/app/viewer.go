package app

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Viewer struct {
	*gtk.Box
}

func NewViewer() *Viewer {
	builder := gtk.NewBuilderFromResource(resourcePrefix + "/viewer.ui")
	container := builder.GetObject("viewer-container").Cast().(*gtk.Box)

	return &Viewer{container}
}

func (w *Viewer) Attach(builder *gtk.Builder) {
	container := builder.GetObject("viewer-root").Cast().(*gtk.ScrolledWindow)
	container.SetChild(w)
}
