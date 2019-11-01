package main

/*
	This piece of code use:
	- gotk3 that is licensed under the ISC License:
	  https://github.com/gotk3/gotk3/blob/master/LICENSE

	- Chroma — A general purpose syntax highlighter in pure Go,
	  under the MIT License: https://github.com/alecthomas/chroma/LICENSE

	- personal libray Copyright ©2018-19 H.F.M - "https://github/hfmrow"
	  This program comes with absolutely no warranty. See the The MIT
	  License (MIT) for details: https://opensource.org/licenses/mit-license.php

	* The purpose of this example is to show you how to highlight source
	  code in your TextView.
*/
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/quick"

	"github.com/gotk3/gotk3/gtk"
)

var (
	textView *gtk.TextView
	win      *gtk.Window
	box      *gtk.Box
	fsBtn    *gtk.FileChooserButton
	err      error
)

func main() {
	gtk.Init(nil)

	win := setupWindow("Go Example HightLight source code in a TextView")

	win.ShowAll()

	// Load and disp source at 1st launch
	loadAndDispSource("syntaxHighlighter.go")

	gtk.Main()
}

// Create and initialize the window
func setupWindow(title string) *gtk.Window {
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create Window:", err)
	}

	win.SetTitle(title)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	width, height := 600, 300
	win.SetDefaultSize(width, height)

	// Box container
	box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	if err != nil {
		log.Fatal("Unable to create Box:", err)
	}

	// To show a text
	textView, err = gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}

	// Allow to scroll the text
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create ScrolledWindow:", err)
	}

	// Give you a chance to load a sourcefile ...
	fsBtn, err = gtk.FileChooserButtonNew("Choose source file", gtk.FILE_CHOOSER_ACTION_OPEN)
	if err != nil {
		log.Fatal("Unable to create FileChooserButton:", err)
	}

	// Set signal callbak to FileChooserButton
	fsBtn.Connect("file-set", filesSelected)

	// Configure TextView
	textView.SetHExpand(true)
	textView.SetVExpand(true)

	// Add widgets to window
	scrolledWindow.Add(textView)
	box.Add(scrolledWindow)
	box.Add(fsBtn)
	win.Add(box)

	return win
}

// filesSelected: callback function for "file-set" signal
func filesSelected(fsb *gtk.FileChooserButton) {
	filename := fsb.GetFilename()
	loadAndDispSource(filename)
}

// loadAndDispSource:
func loadAndDispSource(filename string) {
	var formattedSource string

	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Unable to load file:", err)
	}
	// Get source formatted using pango markup format
	formattedSource, err = ChromaHighlight(string(src))

	// fill TextµBuffer with formatted text
	buff, err := textView.GetBuffer()
	if err != nil {
		log.Fatal("Unable to retrieve TextBuffere:", err)
	}
	// Clean text window before fill it
	buff.Delete(buff.GetStartIter(), buff.GetEndIter())

	// insert markup to the TextBuffer
	buff.InsertMarkup(buff.GetStartIter(), formattedSource)
}

/*
	The following code is part of personal libray (informations above)
*/

// ChromaHighlight: Syntax highlighter using Chroma syntax
// highlighter: "github.com/alecthomas/chroma"
// informations above
func ChromaHighlight(inputString string) (out string, err error) {
	
	buff := new(bytes.Buffer)
	writer := bufio.NewWriter(buff)

	// Registrering pango formatter
	formatters.Register("pango", chroma.FormatterFunc(pangoFormatter))

	// Doing the job (io.Writer, SourceText, language(go), Lexer(pango), style(pygments))
	if err = quick.Highlight(writer, inputString, "go", "pango", "pygments"); err != nil {
		return
	}
	writer.Flush()
	return pangoFinalize(string(buff.Bytes())), err
}

// pangoFormatter: is a part of "ChromaHighlight" function
func pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {

	// Clear the background colour.
	var clearBackground = func(style *chroma.Style) *chroma.Style {
		builder := style.Builder()
		bg := builder.Get(chroma.Background)
		bg.Background = 0
		bg.NoInherit = true
		builder.AddEntry(chroma.Background, bg)
		style, _ = builder.Build()
		return style
	}

	closer, out := "", ""
	style = clearBackground(style)
	for token := it(); token != chroma.EOF; token = it() {
		entry := style.Get(token.Type)
		if !entry.IsZero() {
			closer, out = "", ""
			if entry.Bold == chroma.Yes {
				out += "<b>"
				closer = "</b>" + closer
			}
			if entry.Underline == chroma.Yes {
				out += "<u>"
				closer = "</u>" + closer
			}
			if entry.Italic == chroma.Yes {
				out += "<i>"
				closer = "</i>" + closer
			}
			if entry.Colour.IsSet() {
				out += fmt.Sprintf("<span foreground=\"#%02X%02X%02X\">", entry.Colour.Red(), entry.Colour.Green(), entry.Colour.Blue())
				closer = "</span>" + closer
			}
			if entry.Background.IsSet() {
				out += fmt.Sprintf("<span background=\"#%02X%02X%02X\">", entry.Background.Red(), entry.Background.Green(), entry.Background.Blue())
				closer = "</span>" + closer
			}
			fmt.Fprint(w, out)
		}
		fmt.Fprint(w, pangoPrepare(token.Value))
		if !entry.IsZero() {
			fmt.Fprint(w, closer)
		}
	}
	return nil
}

var pangoEscapeChar = [][]string{{"<", "&lt;", "lOwErThAnTmPrEpLaCeMeNt"}, {"&", "&amp;", "aMpErSaNdTmPrEpLaCeMeNt"}}

// prepare: sanitize input string to safely use with pango
func pangoPrepare(inString string) string {
	inString = strings.ReplaceAll(inString, pangoEscapeChar[1][0], pangoEscapeChar[1][2])
	return strings.ReplaceAll(inString, pangoEscapeChar[0][0], pangoEscapeChar[0][2])
}

// finalize: restore originals characters using markup replacement
func pangoFinalize(inString string) string {
	inString = strings.ReplaceAll(inString, pangoEscapeChar[1][2], pangoEscapeChar[1][1])
	return strings.ReplaceAll(inString, pangoEscapeChar[0][2], pangoEscapeChar[0][1])
}
