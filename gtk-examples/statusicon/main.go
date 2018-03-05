package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var application *gtk.Application
var window *gtk.ApplicationWindow

func buildMenu() *gtk.Menu {
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal("Could not create menu:", err)
	}

	item, err := gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		log.Fatal("Could not create menu item:", err)
	}
	item.Show()
	item.Connect("activate", func() {
		application.Quit()
	})

	menu.Append(item)

	item, err = gtk.MenuItemNewWithLabel("Show App")
	if err != nil {
		log.Fatal("Could not create menu item:", err)
	}
	item.Show()
	item.Connect("activate", func() {
		window.Present()
	})

	menu.Append(item)

	return menu
}

func main() {
	const appID = "com.github.gotk3.gotk3-examples.statusicon"
	var err error
	application, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application:", err)
	}

	application.Connect("activate", func() {

		window, err = gtk.ApplicationWindowNew(application)
		if err != nil {
			log.Fatal("Could not create application window:", err)
		}
		window.SetTitle("StatusIcon Example")
		window.SetPosition(gtk.WIN_POS_MOUSE)
		window.SetDefaultSize(600, 300)

		lbl, err := gtk.LabelNew("Double click the status icon to activate this window. Right click it for a menu.")
		if err != nil {
			log.Fatal("Could not create label:", err)
		}
		window.Add(lbl)

		window.ShowAll()

		menu := buildMenu()

		si, err := gtk.StatusIconNewFromIconName("code")
		if err != nil {
			log.Fatal("Could not create status icon:", err)
		}
		si.SetVisible(true)

		si.Connect("popup-menu", func(statusIcon *gtk.StatusIcon, button, activateTime uint) {
			menu.PopupAtStatusIcon(statusIcon, button, uint32(activateTime))
		})
		si.Connect("activate", func() {
			window.Present()
		})

	})

	os.Exit(application.Run(os.Args))
}
