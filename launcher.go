package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func checkAssets() bool {
	data, err := ioutil.ReadFile(path.Join(cfg.minepath, "assets/indexes", cfg.assetIndex+".json"))

	if err != nil {
		return false
	}

	var assetsJSON assetsFile
	err = json.Unmarshal(data, &assetsJSON)

	if err != nil {
		return false
	}

	for _, v := range assetsJSON.Objects {
		info, err := os.Stat(path.Join(cfg.minepath, "assets/objects", v.Hash[:2], v.Hash))

		if err != nil || info.Size() != v.Size {
			return false
		}
	}

	return true
}

func runmine() error {
	lp, err := getLibsPaths(path.Join(cfg.minepath, "libraries"))

	if err != nil {
		return err
	}

	paths := append(lp, fmt.Sprintf("%s/versions/%s/%s.jar", cfg.minepath, cfg.minever, cfg.minever))
	var strpaths string

	if os := runtime.GOOS; os == "linux" || os == "darwin" {
		strpaths = strings.Join(paths, ":")
	} else {
		strpaths = strings.Join(paths, ";")
	}

	args := []string{
		fmt.Sprintf("-Xmx%sM", cfg.memory),
		fmt.Sprintf("-Djava.library.path=%s/versions/%s/natives/", cfg.minepath, cfg.minever),
		"-cp", strpaths,
		"net.minecraft.client.main.Main",
		"-username", cfg.username,
		"-version", cfg.minever,
		"-assetsDir", cfg.assetsDir,
		"-assetIndex", cfg.assetIndex,
		"-accessToken", cfg.accessToken,
		"-uuid", cfg.uuid,
	}

	cmd := exec.Command("java", args...)
	cmd.Dir = cfg.minepath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
