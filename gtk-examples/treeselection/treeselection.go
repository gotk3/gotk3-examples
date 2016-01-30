package main

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	Window    *gtk.Window
	RootBox   *gtk.Box
	TreeView  *gtk.TreeView
	Entry     *gtk.Entry
	ListStore *gtk.ListStore
)

// Appends single value to the TreeView's model
func AppendToList(value string) {
	ListStore.SetValue(ListStore.Append(), 0, value)
}

// Appends several values to the TreeView's model
func AppendMultipleToList(values ...string) {
	for _, v := range values {
		AppendToList(v)
	}
}

// Handler of "changed" signal of TreeView's selection
func SelectionChanged(s *gtk.TreeSelection) {
	// Returns glib.List of gtk.TreePath pointers
	rows := s.GetSelectedRows(ListStore)
	items := make([]string, 0, rows.Length())

	for l := rows; l != nil; l = l.Next() {
		path := l.Data().(*gtk.TreePath)
		iter, _ := ListStore.GetIter(path)
		value, _ := ListStore.GetValue(iter, 0)
		str, _ := value.GetString()
		items = append(items, str)
	}

	Entry.SetText(fmt.Sprint(items))
}

func main() {
	gtk.Init(&os.Args)

	// Declarations
	Window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	RootBox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	TreeView, _ = gtk.TreeViewNew()
	Entry, _ = gtk.EntryNew()
	ListStore, _ = gtk.ListStoreNew(glib.TYPE_STRING)

	// Window properties
	Window.SetTitle("Products written in Go")
	Window.Connect("destroy", gtk.MainQuit)

	// TreeView properties
	{
		renderer, _ := gtk.CellRendererTextNew()
		column, _ := gtk.TreeViewColumnNewWithAttribute("Value", renderer, "text", 0)
		TreeView.AppendColumn(column)
	}
	TreeView.SetModel(ListStore)

	// TreeView selection properties
	sel, _ := TreeView.GetSelection()
	sel.SetMode(gtk.SELECTION_MULTIPLE)
	sel.Connect("changed", SelectionChanged)

	// Packing
	RootBox.PackStart(TreeView, true, true, 0)
	RootBox.PackStart(Entry, false, false, 0)
	Window.Add(RootBox)

	// Populating list
	// TODO: Add more values to the list
	AppendMultipleToList("Go", "Docker", "CockroachDB")

	Window.ShowAll()
	gtk.Main()
}
