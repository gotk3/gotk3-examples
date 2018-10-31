/*
 * Copyright (c) 2013-2014 Conformal Systems <info@conformal.com>
 *
 * This file originated from: http://opensource.conformal.com/
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

package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func setup_window(title string) *gtk.Window {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

func setup_box(orient gtk.Orientation) *gtk.Box {
	box, err := gtk.BoxNew(orient, 0)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}
	return box
}

func setup_tview() *gtk.TextView {
	tv, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	return tv
}

func setup_btn(label string, onClick func()) *gtk.Button {
	btn, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	btn.Connect("clicked", onClick)
	return btn
}

func get_buffer_from_tview(tv *gtk.TextView) *gtk.TextBuffer {
	buffer, err := tv.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	return buffer
}

func get_text_from_tview(tv *gtk.TextView) string {
	buffer := get_buffer_from_tview(tv)
	start, end := buffer.GetBounds()

	text, err := buffer.GetText(start, end, true)
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	return text
}

func set_text_in_tview(tv *gtk.TextView, text string) {
	buffer := get_buffer_from_tview(tv)
	buffer.SetText(text)
}

// The code before this line is unchanged from the textview example.
// Kept here because we need content to fill our stack.

var (
	winTitle = "stack example"
)

func newBoxText(content string) gtk.IWidget {
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	tv := setup_tview()
	set_text_in_tview(tv, content)
	box.PackStart(tv, true, true, 0)

	btn := setup_btn("Submit", func() {
		text := get_text_from_tview(tv)
		fmt.Println(text)
	})
	box.Add(btn)
	return box
}

func newBoxRadio(btns ...string) gtk.IWidget {
	var (
		// Reference to previous button, so we can add the new one in the same group.
		prev *gtk.RadioButton
		box  = setup_box(gtk.ORIENTATION_VERTICAL)
	)

	for i, txt := range btns {
		radio, err := gtk.RadioButtonNewWithLabelFromWidget(prev, txt)
		if err != nil {
			log.Fatal("Unable to get text:", err)
		}
		box.PackStart(radio, false, false, 0)
		prev = radio

		// We're in a loop, so we need to make a static copy of the index for the callback.
		i := i

		radio.Connect("toggled", func() { fmt.Println(i, radio.GetActive()) })
	}

	return box
}

func newStackFull() gtk.IWidget {
	// get a stack and its switcher.
	stack, err := gtk.StackNew()
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}

	sw, err := gtk.StackSwitcherNew()
	if err != nil {
		log.Fatal("Unable to get text:", err)
	}
	sw.SetStack(stack)

	// Fill the stack with 3 pages.
	boxText1 := newBoxText("Hello there!")
	boxRadio := newBoxRadio("choice 1", "choice 2", "choice 3", "choice 4")
	boxText2 := newBoxText("third page")

	stack.AddTitled(boxText1, "key1", "first page")
	stack.AddTitled(boxRadio, "key2", "second page")
	stack.AddTitled(boxText2, "key3", "third page")

	// You can use icons for a switcher page (the page title will be visible as tooltip).
	stack.ChildSetProperty(boxRadio, "icon-name", "list-add")

	// Pack in a box.
	box := setup_box(gtk.ORIENTATION_VERTICAL)
	box.PackStart(sw, false, false, 0)
	box.PackStart(stack, true, true, 0)
	return box
}

func main() {
	gtk.Init(nil)

	win := setup_window(winTitle)

	box := newStackFull()
	win.Add(box)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
