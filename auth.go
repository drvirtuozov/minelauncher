package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func auth() error {
	res, err := http.PostForm("https://authserver.ely.by/auth/authenticate", url.Values{
		"username":    []string{username.GetText()},
		"password":    []string{password.GetText()},
		"clientToken": []string{cfg.clientToken},
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

	p := jsonRes.SelectedProfile
	p.AccessToken = jsonRes.AccessToken
	err = setProfiles([]profile{p})

	if err != nil {
		return err
	}

	return nil
}
