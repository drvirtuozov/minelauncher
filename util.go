package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

func getLibsPaths(dir string) (paths []string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filepath := path.Join(dir, file.Name())

		if file.IsDir() {
			paths = append(paths, getLibsPaths(filepath)...)
		} else {
			paths = append(paths, filepath)
		}
	}

	return paths
}

func getLauncherProfiles() (lprofiles launcherProfiles, err error) {
	filePath := path.Join(cfg.minepath, "launcher_profiles.json")
	jsonBlob, err := ioutil.ReadFile(filePath)

	if err != nil {
		return lprofiles, err
	}

	err = json.Unmarshal(jsonBlob, &lprofiles)

	if err != nil {
		return lprofiles, err
	}

	return lprofiles, nil
}

func setLauncherProfiles(lprofiles launcherProfiles) error {
	filePath := path.Join(cfg.minepath, "launcher_profiles.json")
	jsonBlob, err := json.Marshal(lprofiles)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, jsonBlob, 0777)

	if err != nil {
		return err
	}

	return nil
}
