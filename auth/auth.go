package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/drvirtuozov/minelauncher/config"
)

const API_ROOT = "https://authserver.ely.by/auth/"

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

func IsAuthenticated() bool {
	if len(config.Runtime.Profiles) == 0 {
		return false
	}

	profile := config.Runtime.Profiles[0]

	if profile.AccessToken != "" && profile.UUID != "" && profile.Name != "" {
		return true
	}

	return false
}

func Authenticate(username, password string) error {
	cfg := config.Runtime
	res, err := http.PostForm(API_ROOT+"authenticate", url.Values{
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

func Logout() error {
	cfg := config.Runtime
	accessToken := cfg.Profiles[0].AccessToken
	cfg.Profiles = []config.Profile{}
	err := config.Set(cfg)

	if err != nil {
		return err
	}

	go http.PostForm(API_ROOT+"invalidate", url.Values{
		"accessToken": []string{accessToken},
		"clientToken": []string{cfg.ClientToken},
	})

	return nil
}
