package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var timeoutContinue = true
var tos glib.SourceHandle

// Create and initialize the window
func setupWindow(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	width, height := 600, 300
	win.SetDefaultSize(width, height)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	btn, _ := gtk.ButtonNew()
	btn.Connect("clicked", buttonClicked)
	btn.SetLabel("Stop timeout")

	box.Add(btn)
	win.Add(box)

	return win
}

func main() {
	gtk.Init(nil)

	win := setupWindow("Go Example Testreport")

	win.ShowAll()

	// Init timeout:
	tos, _ = glib.TimeoutAdd(uint(1000), func() bool {
		fmt.Println("timed out")

		return timeoutContinue
	})

	gtk.Main()
}

func buttonClicked() {
	var methode int

	// Three methods to stop timeout:
	// 1- Destroying
	// 2- Removing
	// 3- Returning false by called function

	// Choose one of them:
	methode = 3

	switch methode {
	case 1: // 1- Destroying
		mcd := glib.MainContextDefault()
		src := mcd.FindSourceById(tos)
		src.Destroy()
		fmt.Printf("Timeout stopped & IsDestroyed: %v\n", src.IsDestroyed())

	case 2: // 2- Removing
		glib.SourceRemove(tos)
		fmt.Printf("Timeout stopped but still referenced\n")

	case 3: // 3- Returning false by called function
		timeoutContinue = !timeoutContinue
		fmt.Printf("Timeout stopped but still in memory and referenced\n")
	}
}
