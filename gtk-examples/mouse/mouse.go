/*
 * Copyright (c) 2020 Sergio Rubio <rubiojr@frameos.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

// Mouse button event handling example.
//
// Demonstrates how to handle left, middle, right mouse button clicks.

package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Mouse Events")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.Add(windowWidget())
	win.ShowAll()

	gtk.Main()
}

func windowWidget() *gtk.Widget {
	lbl, err := gtk.LabelNew("Use the mouse buttons")
	lbl.SetSizeRequest(200, 200)
	if err != nil {
		panic(err)
	}
	evtBox, err := gtk.EventBoxNew()
	if err != nil {
		panic(err)
	}
	evtBox.Add(lbl)
	evtBox.Connect("button-press-event", func(tree *gtk.EventBox, ev *gdk.Event) bool {
		btn := gdk.EventButtonNewFromEvent(ev)
		fmt.Println("button pressed")
		switch btn.Button() {
		case gdk.BUTTON_PRIMARY:
			lbl.SetText("left-click detected!")
			return true
		case gdk.BUTTON_MIDDLE:
			lbl.SetText("middle-click detected!")
			return true
		case gdk.BUTTON_SECONDARY:
			lbl.SetText("right-click detected!")
			return true
		default:
			return false
		}
	})

	return &evtBox.Widget
}
