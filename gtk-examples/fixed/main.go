package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
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
	win.SetTitle("GtkFixed example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	//
	// Important part here:
	//
	// Create a GtkFixed widget to show in the window.
	//
	fixed, err := gtk.FixedNew()
	if err != nil {
		log.Fatal("Unable to create GtkFixed:", err)
	}

	buttonHello, err := gtk.ButtonNewWithLabel("Hello")
	if err != nil {
		log.Fatal("Unable to create Button:", err)
	}

	buttonWorld, err := gtk.ButtonNewWithLabel("World")
	if err != nil {
		log.Fatal("Unable to create Button:", err)
	}

	// Add items to the GtkFixed
	fixed.Put(buttonHello, 100, 200)
	fixed.Put(buttonWorld, 200, 300)

	// Add the label to the window.
	win.Add(fixed)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
