package main

import (
	"fmt"
	"io"
	"strconv"
)

type profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
	Legacy      bool   `json:"legacy"`
}

type launcherProfiles struct {
	ClientToken string    `json:"clientToken"`
	Profiles    []profile `json:"profiles"`
}

type authResponse struct {
	AccessToken     string  `json:"accessToken"`
	ClientToken     string  `json:"clientToken"`
	SelectedProfile profile `json:"selectedProfile"`
}

type authError struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

type assetsFile struct {
	Objects map[string]struct {
		Hash string `json:"hash"`
		Size int64  `json:"size"`
	} `json:"objects"`
}

type passThruReader struct {
	io.Reader
	total  int64
	length int64
}

func (pt *passThruReader) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)

	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)
		fmt.Println("Downloading client update... " + strconv.Itoa(int(percentage)) + "%")
	}

	return n, err
}
