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
	window.SetBorderWidth(10)
	window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.NewVBox(false, 1)
	authBoxAlign := gtk.NewAlignment(0, 1, 0, 0)
	updateBtnAlign := gtk.NewAlignment(0, 0, 0, 0)
	vbox.PackStartDefaults(updateBtnAlign)
	vbox.PackStartDefaults(authBoxAlign)

	authBox := gtk.NewHBox(true, 0)
	authBox.SetSpacing(5)
	authBoxAlign.Add(authBox)
	usernameBox := gtk.NewVBox(true, 0)
	passwordBox := gtk.NewVBox(true, 0)
	usernameLabel := gtk.NewLabel("Username:")
	usernameLabel.SetAlignment(0, 1)
	usernameEntry = gtk.NewEntry()
	passwordLabel := gtk.NewLabel("Password:")
	passwordLabel.SetAlignment(0, 1)
	passwordEntry = gtk.NewEntry()
	passwordEntry.SetVisibility(false)
	authBtn := gtk.NewButtonWithLabel("Authorize via Ely.by")

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
	authBtnAlign := gtk.NewAlignment(0, 1, 0, 0)
	authBtnAlign.Add(authBtn)
	authBox.Add(authBtnAlign)
	updateBtn := gtk.NewButtonWithLabel("Update Client")
	updateBtnAlign.Add(updateBtn)

	updateBtn.Connect("clicked", func() {
		if err := updateClient(); err != nil {
			msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			msg.SetTitle("Update Client Error")
			msg.Response(msg.Destroy)
			msg.Run()
			return
		}

		updateBtn.Destroy()
	})

	go func() {
		if checkClientUpdates() {
			updateBtn.Show()
		}
	}()

	window.Add(vbox)
	window.ShowAll()
	updateBtn.Hide()
	gtk.Main()
}
