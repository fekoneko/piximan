package resources

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed piximan.gresource
var resources []byte

const resourcePrefix = "/com/fekoneko/piximan/"

// Read and register all resources. May be called as soon as possible.
func LoadResources() {
	bytes := glib.NewBytes(resources)
	resource, err := gio.NewResourceFromData(bytes)
	if err != nil {
		panic(err)
	}
	gio.ResourcesRegister(resource)
}

// Create CSS Provider and call StyleContextAddProviderForDisplay().
// Meant to be used once the app is activated.
func LoadCss() {
	cssProvider := gtk.NewCSSProvider()
	cssProvider.LoadFromResource(ResourcePath("window.css"))
	priority := uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	gtk.StyleContextAddProviderForDisplay(gdk.DisplayGetDefault(), cssProvider, priority)
}

// Get builder from resource. Suffix is the part right after the app id without leading slash.
func NewBuilder(suffix string) *gtk.Builder {
	builder := gtk.NewBuilderFromResource(ResourcePath(suffix))
	return builder
}

// Get resource path. Suffix is the part right after the app id without leading slash.
func ResourcePath(suffix string) string {
	return resourcePrefix + suffix
}
