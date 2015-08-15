gotk3 examples
==============

The gotk3 project provides Go bindings for GTK+3 and dependent
projects.  

## Examples for gotk3

## Sample Use

The following example can be found in `gtk-examples/simple/simple.go`.
Usage of additional features is also demonstrated in the
`gtk-examples/` directory.

```Go
package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	l, err := gtk.LabelNew("Hello, gotk3!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Add the label to the window.
	win.Add(l)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run. 
	gtk.Main()
}
```

## Installation

gotk3 currently requires GTK 3.16, GLib 2.36-2.40, and
Cairo 1.10 or 1.12.  A recent Go (1.3 or newer) is also required.

The gtk package requires the cairo, glib, and gdk packages as
dependencies, so only one `go get` is necessary for complete
installation.

The build process uses the tagging scheme gtk_MAJOR_MINOR to specify a
build targeting any particular GTK version (for example, gtk_3_10).
Building with no tags defaults to targeting the latest supported GTK
release (3.16).

To install gotk3 targeting the latest GTK version:

```bash
$ go get github.com/gotk3/gotk3/gtk
```

On MacOS (using homebrew) you would likely specify PKG_CONFIG_PATH as such:
```bash
$ PKG_CONFIG_PATH=/opt/X11/lib/pkgconfig:`brew --prefix gtk+3`/lib/pkgconfig go get -u -v github.com/gotk3/gotk3/gdk
```

```bash
$ sudo apt-get install libgtk-3-dev
$ sudo apt-get install libcairo2-dev
$ sudo apt-get install libglib2.0-dev
```

## License

Package gotk3 is licensed under the liberal ISC License.