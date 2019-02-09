package main

import (
    "log"
    "errors"

    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
)

func main() {
    // Initialize GTK without parsing any command line arguments.
    gtk.Init(nil)

    // Get the GtkBuilder UI definition in the glade file.
    builder, err := gtk.BuilderNewFromFile("ui/example.glade")
    errorCheck(err)

    // Map the handlers to callback functions, and connect the signals 
    // to the Builder. 
    signals := map[string]interface{} {"destroy": on_window_destroy}
    builder.ConnectSignals(signals)

    // Get the object with the id of "window". 
    obj, err := builder.GetObject("window")
    errorCheck(err)

    // Verify that it is a pointer to a Window Widget.
    win, err := isWindow(obj)
    errorCheck(err)
    
    // Show the Widget.
    win.Show()

    // Begin executing the GTK main loop.  This blocks until
    // gtk.MainQuit() is run.
    gtk.Main()
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
    // Make type assertion (as per gtk.go).
    if win, ok := obj.(*gtk.Window); ok {
      return win, nil
    } 
    return nil, errors.New("not a *gtk.Window")
}

func errorCheck(e error) {
    if e != nil {
      // panic for any errors.
      log.Panic(e)
    }
}

func on_window_destroy() {
    log.Println("on_window_destroy")
    // Terminate the GTK main loop.
    gtk.MainQuit()
}