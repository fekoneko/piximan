package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/fekoneko/piximan/internal/collection"
	"github.com/fekoneko/piximan/internal/collection/work"
	"github.com/fekoneko/piximan/internal/downloader/queue"
	"github.com/fekoneko/piximan/internal/logger"
	"github.com/fekoneko/piximan/internal/resources"
)

const openBufferSize = 100

type Window struct {
	*adw.ApplicationWindow
	splitView *adw.NavigationSplitView
	title     *adw.WindowTitle
	explorer  *Explorer
	viewer    *Viewer

	collection *collection.Collection
	works      map[string]*work.StoredWork // TODO: mutex?
	worksMutex *sync.Mutex
}

func NewWindow() *Window {
	builder := resources.NewBuilder("window.ui")

	window := builder.GetObject("window").Cast().(*adw.ApplicationWindow)
	splitView := builder.GetObject("split-view").Cast().(*adw.NavigationSplitView)
	collectionTitle := builder.GetObject("title").Cast().(*adw.WindowTitle)
	works := make(map[string]*work.StoredWork)
	worksMutex := &sync.Mutex{}

	w := &Window{window, splitView, collectionTitle, nil, nil, nil, works, worksMutex}
	w.explorer = NewExplorer(w)
	w.viewer = NewViewer(w)

	w.explorer.Attach(builder)
	w.viewer.Attach(builder)

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
	w.title.SetTitle(title)
	w.Window.SetTitle(fmt.Sprintf("%s - piximan", title))

	if w.collection != nil {
		w.collection.Cancel()
	}

	c := collection.New(path, logger.DefaultLogger)
	works := make(map[string]*work.StoredWork)
	w.collection = c
	w.works = works

	w.explorer.Clear()
	c.Read()

	go func() {
		buffer := make([]string, 0, openBufferSize)

		for work := c.WaitNext(); work != nil; work = c.WaitNext() {
			hash := strconv.Itoa(len(works)) // TODO: different hash?
			w.worksMutex.Lock()
			works[hash] = work
			w.worksMutex.Unlock()

			buffer = append(buffer, hash)
			if len(buffer) >= openBufferSize {
				currentBuffer := buffer
				glib.IdleAdd(func() {
					if !c.Cancelled() {
						w.explorer.Append(currentBuffer...)
					}
				})
				buffer = make([]string, 0, openBufferSize)
			}
		}
		glib.IdleAdd(func() {
			if !c.Cancelled() {
				w.explorer.Append(buffer...)
				logger.Info("%v found in the collection", len(works))
			}
		})

		// TODO: remove
		workWorks := make([]*work.Work, 0, len(works))
		for _, w := range works {
			workWorks = append(workWorks, w.Work)
		}
		l := queue.IgnoreListFromWorks(workWorks)
		logger.Info("ignore list: %v", l.Len())
	}()
}

func (w *Window) ToggleSidebar() {
	collapsed := w.splitView.Collapsed()
	w.splitView.SetShowContent(!collapsed)
	w.splitView.SetCollapsed(!collapsed)
}

func (w *Window) Work(hash string) *work.StoredWork {
	w.worksMutex.Lock()
	defer w.worksMutex.Unlock()
	return w.works[hash]
}
