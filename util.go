package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func getLibsPaths(dir string) (paths []string, err error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filepath := path.Join(dir, file.Name())

		if file.IsDir() {
			filepaths, err := getLibsPaths(filepath)

			if err != nil {
				return nil, err
			}

			paths = append(paths, filepaths...)
		} else {
			paths = append(paths, filepath)
		}
	}

	return paths, nil
}

func getLauncherConfig() (config launcherConfig, err error) {
	filePath := path.Join(minepath, launcher+".json")
	jsonBlob, err := ioutil.ReadFile(filePath)

	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(jsonBlob, &config); err != nil {
		return config, err
	}

	return config, nil
}

func setLauncherConfig(config launcherConfig) error {
	filePath := path.Join(minepath, launcher+".json")
	jsonBlob, err := json.Marshal(config)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filePath, jsonBlob, 0777); err != nil {
		return err
	}

	return nil
}

func isAuthorized() bool {
	if len(cfg.Profiles) == 0 {
		return false
	}

	profile := cfg.Profiles[0]

	if profile.AccessToken != "" && profile.UUID != "" && profile.Name != "" {
		return true
	}

	return false
}

func downloadZip(url string) (filePath string, err error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	file, err := ioutil.TempFile("", launcher+"-update")

	if err != nil {
		return "", err
	}

	defer file.Close()
	reader := &passThruReader{Reader: res.Body, length: res.ContentLength}
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return "", err
	}

	_, err = file.Write(data)

	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func unzip(zipPath, destPath string) error {
	zipReader, err := zip.OpenReader(zipPath)

	if err != nil {
		return err
	}

	if err := os.MkdirAll(destPath, 0777); err != nil {
		return err
	}

	for i, file := range zipReader.File {
		filePath := path.Join(destPath, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		fmt.Printf("Extracting... %d of %d files\n", i+1, len(zipReader.File))
		fileReader, err := file.Open()

		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return nil
}

func copyDir(dirPath, destDir string) error {
	files, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		fromPath := path.Join(dirPath, file.Name())
		toPath := path.Join(destDir, file.Name())

		if file.IsDir() {
			os.MkdirAll(toPath, file.Mode())
			copyDir(fromPath, toPath)
			continue
		}

		bytes, err := ioutil.ReadFile(fromPath)

		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(toPath, bytes, file.Mode()); err != nil {
			return err
		}
	}

	return nil
}

func getCommitFromFilename(filename string) string {
	return strings.TrimSuffix(filename[strings.LastIndex(filename, "-")+1:], path.Ext(filename))
}

func checkClientUpdates() bool {
	res, err := http.Get(cfg.ClientURL)

	if err != nil {
		return false
	}

	defer res.Body.Close()
	header := res.Header.Get("Content-Disposition")
	key := "filename="
	filename := header[strings.Index(header, key)+len(key):]
	commit := getCommitFromFilename(filename)

	if cfg.LastClientCommit != commit {
		return true
	}

	return false
}
