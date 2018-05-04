package launcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/drvirtuozov/minelauncher/config"
	"github.com/drvirtuozov/minelauncher/events"
	"github.com/drvirtuozov/minelauncher/util"
)

type assetsFile struct {
	Objects map[string]struct {
		Hash string `json:"hash"`
		Size int64  `json:"size"`
	} `json:"objects"`
}

func CheckAssets() bool {
	cfg := config.Runtime

	data, err := ioutil.ReadFile(path.Join(cfg.Minepath, "assets/indexes", cfg.AssetIndex+".json"))

	if err != nil {
		return false
	}

	var assetsJSON assetsFile
	err = json.Unmarshal(data, &assetsJSON)

	if err != nil {
		return false
	}

	for _, v := range assetsJSON.Objects {
		info, err := os.Stat(path.Join(cfg.Minepath, "assets/objects", v.Hash[:2], v.Hash))

		if err != nil || info.Size() != v.Size {
			return false
		}
	}

	return true
}

func Runmine() error {
	cfg := config.Runtime

	lp, err := util.GetLibsPaths(path.Join(cfg.Minepath, "libraries"))

	if err != nil {
		return err
	}

	paths := append(lp, fmt.Sprintf("%s/versions/%s/%s.jar", cfg.Minepath, cfg.MinecraftVersion, cfg.MinecraftVersion))
	var strpaths string

	if os := runtime.GOOS; os == "linux" || os == "darwin" {
		strpaths = strings.Join(paths, ":")
	} else {
		strpaths = strings.Join(paths, ";")
	}

	args := []string{
		fmt.Sprintf("-Xmx%dM", cfg.MaxMemory),
		fmt.Sprintf("-Djava.library.path=%s/versions/%s/natives/", cfg.Minepath, cfg.MinecraftVersion),
		"-cp", strpaths,
		"net.minecraft.client.main.Main",
		"-username", cfg.Profiles[0].Name,
		"-version", cfg.MinecraftVersion,
		"-assetsDir", path.Join(cfg.Minepath, "assets"),
		"-assetIndex", cfg.AssetIndex,
		"-accessToken", cfg.Profiles[0].AccessToken,
		"-uuid", cfg.Profiles[0].UUID,
	}

	cmd := exec.Command("java", args...)
	cmd.Dir = cfg.Minepath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func UpdateClient() error {
	cfg := config.Runtime
	zipPath, err := util.DownloadZip(cfg.ClientURL)

	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir("", fmt.Sprintf("%s-client-update", cfg.Launcher))

	if err != nil {
		return err
	}

	if err := util.Unzip(zipPath, tempDir); err != nil {
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

	commit := util.GetCommitFromFilename(dirToRename)
	cfg.LastClientCommit = commit

	events.TaskProgress <- events.ProgressBarFraction{
		Fraction: 1.0,
		Text:     "Copying files...",
	}

	if err := util.CopyDir(dirToRename, cfg.Minepath); err != nil {
		return err
	}

	events.TaskProgress <- events.ProgressBarFraction{
		Fraction: 1.0,
		Text:     "Removing files...",
	}

	config.Set(cfg)
	go os.RemoveAll(tempDir)
	go os.Remove(zipPath)
	return nil
}

func CheckClientUpdates() bool {
	cfg := config.Runtime
	res, err := http.Get(cfg.ClientURL)

	if err != nil {
		return false
	}

	defer res.Body.Close()
	header := res.Header.Get("Content-Disposition")
	key := "filename="
	filename := header[strings.Index(header, key)+len(key):]
	commit := util.GetCommitFromFilename(filename)

	if cfg.LastClientCommit != commit {
		return true
	}

	return false
}
