package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "org.gotk3.example"

func main() {
	app, _ := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)

	//Shows an application as soon as the app starts
	app.Connect("activate", func() {
		notif := glib.NotificationNew("Title")
		notif.SetBody("Text")
		app.SendNotification(appID, notif)
	})

	app.Run(nil)
}
