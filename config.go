package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/joho/godotenv"
)

type config struct {
	launcher    string
	minepath    string
	memory      string
	minever     string
	username    string
	assetsDir   string
	assetIndex  string
	accessToken string
	uuid        string
	clientURL   string
}

var cfg config

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg.launcher = os.Getenv("LAUNCHER")
	cfg.memory = os.Getenv("MEMORY")
	cfg.minever = os.Getenv("MINEVER")
	cfg.username = os.Getenv("USERNAME")
	cfg.assetIndex = os.Getenv("ASSET_INDEX")
	cfg.accessToken = os.Getenv("ACCESS_TOKEN")
	cfg.uuid = os.Getenv("UUID")
	cfg.clientURL = os.Getenv("CLIENT_URL")
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	cfg.minepath = fmt.Sprintf("%s/.%s", usr.HomeDir, cfg.launcher)
	cfg.assetsDir = cfg.minepath + "/assets/"
}
