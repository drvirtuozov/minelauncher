package main

import (
	"github.com/mattn/go-gtk/gtk"
)

var username *gtk.Entry
var password *gtk.Entry

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(cfg.launcher)
	window.SetSizeRequest(800, 400)
	window.SetResizable(false)
	window.Connect("destroy", gtk.MainQuit)
	fixed := gtk.NewFixed()
	authBox := gtk.NewHBox(true, 0)
	usernameBox := gtk.NewVBox(true, 0)
	passwordBox := gtk.NewVBox(true, 0)
	usernameLabel := gtk.NewLabel("Username:")
	username = gtk.NewEntry()
	passwordLabel := gtk.NewLabel("Password:")
	password = gtk.NewEntry()
	password.SetVisibility(false)
	authBtn := gtk.NewButton()
	authBtn.SetLabel("Authorize via Ely.by")
	authBtn.Connect("clicked", func() {
		if err := auth(); err != nil {
			msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			msg.SetTitle("Authorization Error")
			msg.Response(msg.Destroy)
			msg.Run()
		}
	})
	usernameBox.Add(usernameLabel)
	usernameBox.Add(username)
	passwordBox.Add(passwordLabel)
	passwordBox.Add(password)
	authBox.Add(usernameBox)
	authBox.Add(passwordBox)
	authBox.Add(authBtn)
	updateBtn := gtk.NewButton()
	updateBtn.SetLabel("Update Client")
	updateBtn.Connect("clicked", func() {
		if err := updateClient(); err != nil {
			msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			msg.SetTitle("Client Updating Error")
			msg.Response(msg.Destroy)
			msg.Run()
		}
	})
	go func() {
		if checkClientUpdates() {
			updateBtn.Show()
		}
	}()
	fixed.Put(authBox, 10, 340)
	fixed.Put(updateBtn, 10, 10)
	window.Add(fixed)
	window.ShowAll()
	updateBtn.Hide()
	gtk.Main()
}
