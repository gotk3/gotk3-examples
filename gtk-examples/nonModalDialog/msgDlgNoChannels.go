package main

import (
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "log"
	"time"
)

func main() {
    gtk.Init(nil)
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetDefaultSize(500, 200)
    win.SetTitle("MessageDialog Example")
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
    msgLbl, err := gtk.LabelNew("")
	if err != nil {
        log.Fatal("Unable to create msgLbl:", err)
    }
	submitBtn, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Unable to create submitBtn:", err)
	}
	submitBtn.SetLabel("Submit")
	submitBtn.Connect("clicked", func() {
		go func() {
			workingDlg := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_NONE, "%s", "Working...")
			_, err := glib.IdleAdd(workingDlg.Show)
			if err != nil {
				log.Fatal("Unable to execute IdleAdd():", err)
			}
			msg := getMsg()
			time.Sleep(1 * time.Second) // This is here only for testing
			_, err2 := glib.IdleAdd(msgLbl.SetText, msg)
			if err2 != nil {
				log.Fatal("Unable to execute IdleAdd():", err2)
			}
			_, err3 := glib.IdleAdd(workingDlg.Destroy)
			if err3 != nil {
				log.Fatal("Unable to execute IdleAdd():", err3)
			}
		}()
	})
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create mainBox:", err)
	}
	mainBox.PackStart(msgLbl, false, false, 0)
	mainBox.PackStart(submitBtn, false, false, 0)
	win.Add(mainBox)
    win.ShowAll()
    gtk.Main()
}

func getMsg() string {
	now := time.Now()
	hms := now.Format("15:04:05")
	return "Time is " + hms
}

