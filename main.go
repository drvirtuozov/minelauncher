package main

import (
	"fmt"
	"log"
	"os/user"
	"path"

	"github.com/drvirtuozov/minelauncher/auth"
	"github.com/drvirtuozov/minelauncher/config"
	"github.com/drvirtuozov/minelauncher/events"
	"github.com/drvirtuozov/minelauncher/launcher"
	"github.com/google/uuid"
	"github.com/mattn/go-gtk/gtk"
)

var lname string
var mversion string
var assetIndex string
var clientURL string

func init() {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	config.Runtime.Launcher = lname
	config.Runtime.Minepath = path.Join(usr.HomeDir, "."+lname)

	cfg, err := config.Get()

	if err != nil {
		log.Fatal(err)
	}

	if cfg.Launcher == "" {
		cfg.Launcher = lname
	}

	if cfg.MinecraftVersion == "" {
		cfg.MinecraftVersion = mversion
	}

	if cfg.AssetIndex == "" {
		cfg.AssetIndex = assetIndex
	}

	if cfg.ClientURL == "" {
		cfg.ClientURL = clientURL
	}

	if cfg.Minepath == "" {
		cfg.Minepath = path.Join(usr.HomeDir, "."+cfg.Launcher)
	}

	if cfg.MaxMemory == 0 {
		cfg.MaxMemory = 1024
	}

	if cfg.ClientToken == "" {
		id := uuid.New()
		cfg.ClientToken = id.String()
	}

	err = config.Set(cfg)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle(lname)
	window.SetSizeRequest(800, 400)
	window.SetResizable(false)
	window.SetBorderWidth(10)
	window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.NewVBox(false, 0)
	hbox := gtk.NewHBox(false, 0)
	authBoxAlign := gtk.NewAlignment(0, 1, 0, 0)
	logoutBoxAlign := gtk.NewAlignment(0, 1, 0, 0)
	updateBtnAlign := gtk.NewAlignment(0, 0, 0, 0)
	playBtnAlign := gtk.NewAlignment(0.5, 1, 0, 0.5)
	progressBarAlign := gtk.NewAlignment(1, 1, 0, 0)

	hbox.Add(logoutBoxAlign)
	hbox.Add(authBoxAlign)
	hbox.Add(progressBarAlign)
	vbox.PackStart(updateBtnAlign, true, true, 0)
	vbox.PackStart(playBtnAlign, true, true, 100)
	vbox.PackStart(hbox, true, true, 0)

	playBtn := gtk.NewButtonWithLabel("Enter the Game")
	playBtnAlign.Add(playBtn)

	progressBar := gtk.NewProgressBar()
	progressBar.SetSizeRequest(350, 26)
	progressBarAlign.Add(progressBar)

	authBox := gtk.NewHBox(true, 0)
	authBox.SetSpacing(5)
	authBoxAlign.Add(authBox)
	usernameBox := gtk.NewVBox(true, 0)
	passwordBox := gtk.NewVBox(true, 0)
	usernameLabel := gtk.NewLabel("Username:")
	usernameLabel.SetAlignment(0, 1)
	usernameEntry := gtk.NewEntry()
	passwordLabel := gtk.NewLabel("Password:")
	passwordLabel.SetAlignment(0, 1)
	passwordEntry := gtk.NewEntry()
	passwordEntry.SetVisibility(false)
	authBtn := gtk.NewButtonWithLabel("Authenticate via Ely.by")

	logoutBox := gtk.NewHBox(true, 0)
	logoutBox.SetSpacing(5)
	logoutBoxAlign.Add(logoutBox)

	logoutBoxLblAlign := gtk.NewAlignment(0, 1, 0, 0)
	logoutBoxLbl := gtk.NewLabel("")
	logoutBoxLbl.SetPadding(0, 5)
	logoutBoxLblAlign.Add(logoutBoxLbl)
	logoutBox.Add(logoutBoxLblAlign)

	logoutBtn := gtk.NewButtonWithLabel("Log Out")
	logoutBtnAlign := gtk.NewAlignment(0, 1, 0, 0)
	logoutBtnAlign.Add(logoutBtn)
	logoutBox.Add(logoutBtnAlign)

	logoutBtn.Connect("clicked", func() {
		if err := auth.Logout(); err != nil {
			log.Println(err)
			return
		}

		logoutBoxAlign.Hide()
		authBoxAlign.Show()
		playBtn.SetSensitive(false)
	})

	authBtn.Connect("clicked", func() {
		username := usernameEntry.GetText()
		password := passwordEntry.GetText()

		if err := auth.Authenticate(username, password); err != nil {
			msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
			msg.SetTitle("Authorization Error")
			msg.Response(msg.Destroy)
			msg.Run()
			return
		}

		authBoxAlign.Hide()
		logoutBoxAlign.Show()
		playBtn.SetSensitive(true)
		logoutBoxLbl.SetText(fmt.Sprintf("Hello, %s", username))
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
		progressBar.Show()
		events.TaskProgress <- events.ProgressBarFraction{
			Text: "Updating client...",
		}
		updateBtn.SetSensitive(false)

		go func() {
			if err := launcher.UpdateClient(); err != nil {
				msg := gtk.NewMessageDialog(window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, err.Error())
				msg.SetTitle("Update Client Error")
				msg.Response(msg.Destroy)
				msg.Run()
				return
			}

			updateBtn.Hide()
			progressBar.Hide()
		}()
	})

	go func() {
		if launcher.CheckClientUpdates() {
			updateBtn.Show()
		}
	}()

	go func() {
		for fraction := range events.TaskProgress {
			progressBar.SetFraction(fraction.Fraction)
			progressBar.SetText(fraction.Text)
		}
	}()

	window.Add(vbox)
	window.ShowAll()
	updateBtn.Hide()
	progressBar.Hide()

	if auth.IsAuthenticated() {
		authBoxAlign.Hide()
		logoutBoxLbl.SetText(fmt.Sprintf("Hello, %s", config.Runtime.Profiles[0].Name))
	} else {
		playBtn.SetSensitive(false)
		logoutBoxAlign.Hide()
	}

	gtk.Main()
}
