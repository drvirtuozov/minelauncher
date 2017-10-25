package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Legacy      bool   `json:"legacy"`
}

type launcherProfiles struct {
	ClientToken string  `json:"clientToken"`
	Profiles    profile `json:"profiles"`
}

type authResponse struct {
	AccessToken     string  `json:"accessToken"`
	ClientToken     string  `json:"clientToken"`
	SelectedProfile profile `json:"selectedProfile"`
}

func auth() error {
	var lprofiles launcherProfiles
	filePath := path.Join(cfg.minepath, "launcher_profiles.json")
	jsonBlob, err := ioutil.ReadFile(filePath)

	if err != nil {
		panic(err)
		return err
	}

	err = json.Unmarshal(jsonBlob, &lprofiles)

	if err != nil {
		panic(err)
		return err
	}

	res, err := http.PostForm("https://authserver.ely.by/auth/authenticate", url.Values{
		"username":    []string{username.GetTooltipText()},
		"password":    []string{password.GetTooltipText()},
		"clientToken": []string{lprofiles.ClientToken},
	})

	if err != nil {
		panic(err)
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	var jsonRes authResponse
	err = json.NewDecoder(res.Body).Decode(&jsonRes)

	if err != nil {
		panic(err)
		return err
	}

	lprofiles.Profiles = jsonRes.SelectedProfile
	lprofiles.Profiles.AccessToken = jsonRes.AccessToken
	jsonBlob, err = json.Marshal(lprofiles)

	if err != nil {
		panic(err)
		return err
	}

	err = ioutil.WriteFile(filePath, jsonBlob, 0777)

	if err != nil {
		panic(err)
		return err
	}

	return nil
}
