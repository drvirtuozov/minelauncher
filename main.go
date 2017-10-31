package main

import (
	"os/user"
	"path"

	"github.com/google/uuid"
	"github.com/mattn/go-gtk/gtk"
)

var launcher string
var minepath string
var cfg launcherConfig
var usernameEntry *gtk.Entry
var passwordEntry *gtk.Entry

func init() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	minepath = path.Join(usr.HomeDir, "."+launcher)
	cfg, _ = getLauncherConfig()

	if cfg.MaxMemory == 0 {
		cfg.MaxMemory = 1024
	}

	if cfg.ClientToken == "" {
		id := uuid.New()
		cfg.ClientToken = id.String()
	}

	setLauncherConfig(cfg)
}

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(launcher)
	window.SetSizeRequest(800, 400)
	window.SetResizable(false)
	window.Connect("destroy", gtk.MainQuit)
	fixed := gtk.NewFixed()
	authBox := gtk.NewHBox(true, 0)
	usernameBox := gtk.NewVBox(true, 0)
	passwordBox := gtk.NewVBox(true, 0)
	usernameLabel := gtk.NewLabel("Username:")
	usernameEntry = gtk.NewEntry()
	passwordLabel := gtk.NewLabel("Password:")
	passwordEntry = gtk.NewEntry()
	passwordEntry.SetVisibility(false)
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
	usernameBox.Add(usernameEntry)
	passwordBox.Add(passwordLabel)
	passwordBox.Add(passwordEntry)
	authBox.Add(usernameBox)
	authBox.Add(passwordBox)
	authBox.Add(authBtn)
	updateBtn := gtk.NewButton()
	updateBtn.SetLabel("Update Client")
	updateBtn.Connect("clicked", func() {
		if err := updateClient(); err != nil {
			msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			msg.SetTitle("Update Client Error")
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
