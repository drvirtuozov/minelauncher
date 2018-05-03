package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/drvirtuozov/minelauncher/config"
)

type responseProfile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Legacy      bool   `json:"legacy"`
}

type authResponse struct {
	AccessToken     string          `json:"accessToken"`
	ClientToken     string          `json:"clientToken"`
	SelectedProfile responseProfile `json:"selectedProfile"`
}

type authError struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func isAuthorized() bool {
	cfg, err := config.Get()

	if err != nil {
		return false
	}

	if len(cfg.Profiles) == 0 {
		return false
	}

	profile := cfg.Profiles[0]

	if profile.AccessToken != "" && profile.UUID != "" && profile.Name != "" {
		return true
	}

	return false
}

func Authenticate(username, password string) error {
	cfg, err := config.Get()

	if err != nil {
		return err
	}

	res, err := http.PostForm("https://authserver.ely.by/auth/authenticate", url.Values{
		"username":    []string{username},
		"password":    []string{password},
		"clientToken": []string{cfg.ClientToken},
	})

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		var jsonRes authError
		err = json.NewDecoder(res.Body).Decode(&jsonRes)
		return errors.New(jsonRes.ErrorMessage)
	}

	var jsonRes authResponse
	err = json.NewDecoder(res.Body).Decode(&jsonRes)

	if err != nil {
		return err
	}

	var profile config.Profile
	profile.UUID = jsonRes.SelectedProfile.ID
	profile.Name = jsonRes.SelectedProfile.Name
	profile.AccessToken = jsonRes.AccessToken
	cfg.Profiles = []config.Profile{profile}
	err = config.Set(cfg)

	if err != nil {
		return err
	}

	return nil
}