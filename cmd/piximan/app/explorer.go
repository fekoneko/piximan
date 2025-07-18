package app

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Explorer struct {
	*gtk.ListView
}

func NewExplorer() *Explorer {
	builder := gtk.NewBuilderFromResource(resourcePrefix + "/explorer.ui")
	list := builder.GetObject("explorer-list").Cast().(*gtk.ListView)

	model := gio.NewListStore(glib.TypeObject)
	factory := gtk.NewBuilderListItemFactoryFromResource(
		builder.Scope(), resourcePrefix+"/card.ui",
	)

	list.ConnectActivate(func(position uint) {
		fmt.Println(position)
	})

	list.SetModel(gtk.NewSingleSelection(model))
	list.SetFactory(&factory.ListItemFactory)

	objects := make([]*glib.Object, 0, 10000)
	for i := range 1000000 {
		item := gtk.NewStringObject(fmt.Sprintf("Work %v", i))
		objects = append(objects, item.Object)
	}

	model.Splice(0, 0, objects)

	return &Explorer{list}
}

func (w *Explorer) Attach(builder *gtk.Builder) {
	container := builder.GetObject("explorer-root").Cast().(*gtk.ScrolledWindow)
	container.SetChild(w)
}
