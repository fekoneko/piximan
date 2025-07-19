package ui

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/resources"
)

type Explorer struct {
	*gtk.ListView
}

func NewExplorer() *Explorer {
	builder := resources.NewBuilder("explorer.ui")
	list := builder.GetObject("explorer-list").Cast().(*gtk.ListView)

	model := gio.NewListStore(glib.TypeObject)
	factory := gtk.NewSignalListItemFactory()

	// TODO: Ideally figure out how to do subclassing and use BuilderListItemFactory
	//       with card template. Store all needed fields in the model and bind them
	//       in the card template: `label: bind template.item as <$CardModel>.title.string`

	factory.ConnectSetup(func(object *glib.Object) {
		item := object.Cast().(*gtk.ListItem)
		card := NewCard()
		item.SetChild(card)
	})

	factory.ConnectBind(func(object *glib.Object) {
		item := object.Cast().(*gtk.ListItem)
		card := CardFromBox(item.Child().(*gtk.Box))
		index := item.Item().Cast().(*gtk.StringObject).String()
		title := fmt.Sprintf("Test Work #%v", index)
		card.Patch(&work.Work{Title: &title})
	})

	list.ConnectActivate(func(position uint) {
		fmt.Println(position)
	})

	list.SetModel(gtk.NewSingleSelection(model))
	list.SetFactory(&factory.ListItemFactory)

	objects := make([]*glib.Object, 0, 10000)
	for i := range 1000000 {
		item := gtk.NewStringObject(strconv.Itoa(i))
		objects = append(objects, item.Object)
	}

	model.Splice(0, 0, objects)

	return &Explorer{list}
}

func (w *Explorer) Attach(builder *gtk.Builder) {
	container := builder.GetObject("explorer-root").Cast().(*gtk.ScrolledWindow)
	container.SetChild(w)
}
