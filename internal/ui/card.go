package ui

import (
	"strings"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/resources"
)

type Card struct {
	*gtk.Box
	thumbnail *gtk.Picture
	title     *gtk.Label
	userName  *gtk.Label
	tags      *gtk.Label
}

func NewCard() *Card {
	builder := resources.NewBuilder("card.ui")
	box := builder.GetObject("card-container").Cast().(*gtk.Box)

	return CardFromBox(box)
}

func CardFromBox(box *gtk.Box) *Card {
	thumbnail := box.FirstChild().(*gtk.Picture)
	infoContainer := thumbnail.NextSibling().(*gtk.Box)
	title := infoContainer.FirstChild().(*gtk.Label)
	username := title.NextSibling().(*gtk.Label)
	tags := username.NextSibling().(*gtk.Label)

	return &Card{box, thumbnail, title, username, tags}
}

func (c *Card) Patch(w *work.StoredWork) {
	// if w != nil && len(w.AssetNames) > 0 {
	// 	thumbnailPath := filepath.Join(w.Path, w.AssetNames[0])
	// 	c.thumbnail.SetFilename(thumbnailPath) // TODO: performant thumbnail decoding
	// } else {
	// 	c.thumbnail.SetFile(nil) // TODO: placeholder
	// }
	if w != nil && w.Title != nil {
		c.title.SetText(*w.Title)
	} else {
		c.title.SetText("Unknown")
	}
	if w != nil && w.UserName != nil {
		c.userName.SetText(*w.UserName)
	} else {
		c.userName.SetText("Unknown")
	}
	if w != nil && w.Tags != nil {
		c.tags.SetText(strings.Join(*w.Tags, ", "))
	} else {
		c.tags.SetText("Unknown")
	}
}
