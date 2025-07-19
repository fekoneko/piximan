package ui

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/resources"
)

type Window struct {
	*adw.ApplicationWindow
	splitView       *adw.NavigationSplitView
	collectionTitle *adw.WindowTitle
	explorer        *Explorer
	viewer          *Viewer
}

func NewWindow() *Window {
	builder := resources.NewBuilder("window.ui")

	window := builder.GetObject("window").Cast().(*adw.ApplicationWindow)
	splitView := builder.GetObject("split-view").Cast().(*adw.NavigationSplitView)
	collectionTitle := builder.GetObject("collection-title").Cast().(*adw.WindowTitle)
	explorer := NewExplorer()
	viewer := NewViewer()
	w := &Window{window, splitView, collectionTitle, explorer, viewer}

	explorer.Attach(builder)
	viewer.Attach(builder)

	openButton := builder.GetObject("open-button").Cast().(*gtk.Button)
	openButton.ConnectClicked(w.AskOpenCollection)

	sidebarToggle := builder.GetObject("sidebar-toggle").Cast().(*gtk.Button)
	sidebarToggle.ConnectClicked(w.ToggleSidebar)

	return w
}

func (w *Window) Attach(app *adw.Application) {
	app.AddWindow(&w.Window)
}

func (w *Window) AskOpenCollection() {
	dialog := gtk.NewFileDialog()
	dialog.SetTitle("Open collection")

	dialog.SelectFolder(context.Background(), &w.Window, func(result gio.AsyncResulter) {
		if file, err := dialog.SelectFolderFinish(result); err == nil {
			w.OpenCollection(file.Path())
		}
	})
}

func (w *Window) OpenCollection(path string) {
	title := filepath.Base(path)
	w.collectionTitle.SetTitle(title)
	w.Window.SetTitle(fmt.Sprintf("%s - piximan", title))
}

func (w *Window) ToggleSidebar() {
	collapsed := w.splitView.Collapsed()
	w.splitView.SetShowContent(!collapsed)
	w.splitView.SetCollapsed(!collapsed)
}
