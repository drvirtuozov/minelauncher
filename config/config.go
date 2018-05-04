package config

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type Profile struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

type Config struct {
	Launcher         string    `json:"launcher"`
	Minepath         string    `json:"minepath"`
	MinecraftVersion string    `json:"minecraft_version"`
	AssetIndex       string    `json:"asset_index"`
	ClientURL        string    `json:"client_url"`
	ClientToken      string    `json:"client_token"`
	MaxMemory        int       `json:"max_memory"`
	Profiles         []Profile `json:"profiles"`
	LastClientCommit string    `json:"last_client_commit"`
}

var Runtime Config

func Get() (cfg Config, err error) {
	filepath := path.Join(Runtime.Minepath, Runtime.Launcher+".json")
	bjson, err := ioutil.ReadFile(filepath)

	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(bjson, &cfg); err != nil {
		return cfg, err
	}

	Runtime = cfg
	return cfg, nil
}

func Set(cfg Config) error {
	filepath := path.Join(Runtime.Minepath, Runtime.Launcher+".json")
	bjson, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath, bjson, 0777); err != nil {
		return err
	}

	Runtime = cfg
	return nil
}
