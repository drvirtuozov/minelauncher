package main

import (
	"fmt"
	"os/user"
)

var (
	launcher   string
	minever    string
	assetIndex string
	clientURL  string
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
	cfg.launcher = launcher
	cfg.minever = minever
	cfg.assetIndex = assetIndex
	cfg.clientURL = clientURL
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	cfg.minepath = fmt.Sprintf("%s/.%s", usr.HomeDir, cfg.launcher)
	cfg.assetsDir = cfg.minepath + "/assets/"
	fmt.Println(cfg)
}
