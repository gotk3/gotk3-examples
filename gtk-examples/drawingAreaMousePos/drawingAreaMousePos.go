// drawingAreaMousePos.go

// from: github.com/hfmrow/


package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var lblX, lblY *gtk.Label

// Create and initialize the window
func setupWindow(title string) *gtk.Window {
	Width, Height := 640, 480

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(Width, Height)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	lblX, _ = gtk.LabelNew("X")
	lblY, _ = gtk.LabelNew("Y")

	da, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatalf("DrawingAreaNew: %s\n", err.Error())
	}
	
	// Adding event we want to recieve (gdk.EVENT_MOTION_NOTIFY)
	// There is a problem using EVENT_MOTION_NOTIFY, it is not set to 4 as it must do
	// So i use simply 4 that correspond to the desired event
	da.AddEvents(4) // accept only integer ...

	// Connect to event we have set previously
	da.Connect("event", daEvent)

	// Setting parameter for drawing area
	da.SetHAlign(gtk.ALIGN_FILL)
	da.SetVAlign(gtk.ALIGN_FILL)
	da.SetHExpand(true)
	da.SetVExpand(true)
	da.SetSizeRequest(Width, Height)

	// Put da in the box
	box.Add(da)
	box.Add(lblX)
	box.Add(lblY)
	
	// And the box into the window ...
	win.Add(box)

	return win
}

// Handling event ...
func daEvent(da *gtk.DrawingArea, event *gdk.Event) bool {
	eventMotion := gdk.EventMotionNewFromEvent(event)
	x, y := eventMotion.MotionVal()
	lblX.SetLabel(fmt.Sprintf("%d", int(x)))
	lblY.SetLabel(fmt.Sprintf("%d", int(y)))
	return false
}

func main() {
	gtk.Init(nil)

	win := setupWindow("Go Example Drawing area mouse event")

	win.ShowAll()

	gtk.Main()
}
