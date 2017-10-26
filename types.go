package main

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
