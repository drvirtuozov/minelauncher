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
	data, err := ioutil.ReadFile(path.Join(minepath, "assets/indexes", cfg.AssetIndex+".json"))

	if err != nil {
		return false
	}

	var assetsJSON assetsFile
	err = json.Unmarshal(data, &assetsJSON)

	if err != nil {
		return false
	}

	for _, v := range assetsJSON.Objects {
		info, err := os.Stat(path.Join(minepath, "assets/objects", v.Hash[:2], v.Hash))

		if err != nil || info.Size() != v.Size {
			return false
		}
	}

	return true
}

func runmine() error {
	lp, err := getLibsPaths(path.Join(minepath, "libraries"))

	if err != nil {
		return err
	}

	paths := append(lp, fmt.Sprintf("%s/versions/%s/%s.jar", minepath, cfg.MinecraftVersion, cfg.MinecraftVersion))
	var strpaths string

	if os := runtime.GOOS; os == "linux" || os == "darwin" {
		strpaths = strings.Join(paths, ":")
	} else {
		strpaths = strings.Join(paths, ";")
	}

	args := []string{
		fmt.Sprintf("-Xmx%dM", cfg.MaxMemory),
		fmt.Sprintf("-Djava.library.path=%s/versions/%s/natives/", minepath, cfg.MinecraftVersion),
		"-cp", strpaths,
		"net.minecraft.client.main.Main",
		"-username", cfg.Profiles[0].Name,
		"-version", cfg.MinecraftVersion,
		"-assetsDir", path.Join(minepath, "assets"),
		"-assetIndex", cfg.AssetIndex,
		"-accessToken", cfg.Profiles[0].AccessToken,
		"-uuid", cfg.Profiles[0].UUID,
	}

	cmd := exec.Command("java", args...)
	cmd.Dir = minepath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func updateClient() error {
	zipPath, err := downloadZip(cfg.ClientURL)

	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir("", fmt.Sprintf("%s-client-update", launcher))

	if err != nil {
		return err
	}

	if err := unzip(zipPath, tempDir); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tempDir)

	if err != nil {
		return err
	}

	var dirToRename string

	if len(files) == 1 && files[0].IsDir() {
		dirToRename = path.Join(tempDir, files[0].Name())
	} else {
		dirToRename = tempDir
	}

	commit := getCommitFromFilename(dirToRename)
	cfg.LastClientCommit = commit

	if err := copyDir(dirToRename, minepath); err != nil {
		return err
	}

	setLauncherConfig(cfg)
	go os.RemoveAll(tempDir)
	go os.Remove(zipPath)
	return nil
}
