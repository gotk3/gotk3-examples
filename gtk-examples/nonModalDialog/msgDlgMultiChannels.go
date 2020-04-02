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
			showWorkDlgChan := make(chan bool)
			go showWorkDlg(workDlg, showWorkDlgChan)
			<-showWorkDlgChan

			time.Sleep(1 * time.Second) // This is here for testing only

			msgChan := make(chan string)
			go getMsg(msgChan)
			msg := <-msgChan

			msgLblChan := make(chan bool)
			go setMsgLbl(msgLbl, msg, msgLblChan)
			<-msgLblChan

			destroyWorkDlgChan := make(chan bool, 1)
			go destroyWorkDlg(workDlg, destroyWorkDlgChan)
			<-destroyWorkDlgChan
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

func showWorkDlg(dlg *gtk.MessageDialog, done chan bool) {
	_, err := glib.IdleAdd(dlg.Show)
	if err != nil {
		log.Fatal("Unable to execute IdleAdd():", err)
	}
	done <- true
}

func getMsg(msgChan chan string) {
	now := time.Now()
	hms := now.Format("15:04:05")
	msgChan <- hms
}

func setMsgLbl(msgLbl *gtk.Label, msg string, msgLblChan chan bool) {
	_, err := glib.IdleAdd(msgLbl.SetText, msg)
	if err != nil {
		log.Fatal("Unable to execute IdleAdd():", err)
	}
	msgLblChan <- true
}

func destroyWorkDlg(dlg *gtk.MessageDialog, done chan bool) {
	_, err := glib.IdleAdd(dlg.Destroy)
	if err != nil {
		log.Fatal("Unable to execute IdleAdd():", err)
	}
	done <- true
}

