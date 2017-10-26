package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
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

func getProfiles() (profiles []profile, err error) {
	filePath := path.Join(cfg.minepath, "launcher_profiles.json")
	jsonBlob, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	var lprofiles launcherProfiles
	err = json.Unmarshal(jsonBlob, &lprofiles)

	if err != nil {
		return lprofiles.Profiles, err
	}

	return lprofiles.Profiles, nil
}

func setProfiles(profiles []profile) error {
	var lprofiles launcherProfiles
	lprofiles.ClientToken = cfg.clientToken
	lprofiles.Profiles = profiles
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

func isAuthorized() bool {
	profiles, err := getProfiles()

	if len(profiles) == 0 || err != nil {
		return false
	}

	profile := profiles[0]

	if profile.AccessToken != "" && profile.ID != "" && profile.Name != "" {
		return true
	}

	return false
}
