package main

import (
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

	cfg, _ := config.Get()

	cfg.Launcher = lname
	cfg.MinecraftVersion = mversion
	cfg.AssetIndex = assetIndex
	cfg.ClientURL = clientURL
	cfg.Minepath = path.Join(usr.HomeDir, "."+cfg.Launcher)

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
	updateBtnAlign := gtk.NewAlignment(0, 0, 0, 0)
	playBtnAlign := gtk.NewAlignment(0.5, 1, 0, 0.5)
	progressBarAlign := gtk.NewAlignment(1, 1, 0, 0)
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
	authBtn := gtk.NewButtonWithLabel("Authorize via Ely.by")

	authBtn.Connect("clicked", func() {
		username := usernameEntry.GetText()
		password := passwordEntry.GetText()

		if err := auth.Authenticate(username, password); err != nil {
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
	gtk.Main()
}
