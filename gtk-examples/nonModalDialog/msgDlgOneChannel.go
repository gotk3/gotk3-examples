package main

import (
    "log"
	"time"

    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
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
		workDlg := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_NONE, "%s", "Working...")
		go func() {
			_, err := glib.IdleAdd(workDlg.Show)
			if err != nil {
				log.Fatal("Unable to execute IdleAdd():", err)
			}
			time.Sleep(1 * time.Second) // This is here for testing only
			msgChan := make(chan string)
			go getMsg(msgChan)
			msg := <-msgChan
			_, err2 := glib.IdleAdd(msgLbl.SetText, msg)
			if err2 != nil {
				log.Fatal("Unable to execute IdleAdd():", err2)
			}
			_, err3 := glib.IdleAdd(workDlg.Destroy)
			if err != nil {
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

func getMsg(msgChan chan string) {
	now := time.Now()
	hms := now.Format("15:04:05")
	msgChan <- hms
}

