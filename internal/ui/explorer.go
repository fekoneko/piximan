package ui

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/resources"
)

type Explorer struct {
	*gtk.ListView
	window *Window
	model  *gio.ListStore
}

func NewExplorer(window *Window) *Explorer {
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
		hash := item.Item().Cast().(*gtk.StringObject).String()
		work := window.Work(hash)
		card.Patch(work)
	})

	list.ConnectActivate(func(position uint) {
		fmt.Println(position)
	})

	list.SetModel(gtk.NewSingleSelection(model))
	list.SetFactory(&factory.ListItemFactory)

	return &Explorer{list, window, model}
}

func (e *Explorer) Attach(builder *gtk.Builder) {
	container := builder.GetObject("explorer-root").Cast().(*gtk.ScrolledWindow)
	container.SetChild(e)
}

// Append work to a list. Provided hash will be used to get the work
// from the parent window when list item is activated.
func (e *Explorer) Append(hashes ...string) {
	if len(hashes) == 1 {
		object := gtk.NewStringObject(hashes[0]).Object
		e.model.Append(object)
	} else if len(hashes) > 1 {
		objects := make([]*glib.Object, 0, len(hashes))
		for _, hash := range hashes {
			object := gtk.NewStringObject(hash).Object
			objects = append(objects, object)
		}
		e.model.Splice(e.model.NItems(), 0, objects)
	}
}

func (e *Explorer) Clear() {
	e.model.RemoveAll()
}
