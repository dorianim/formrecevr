package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type TokenValidateRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIP string `json:"remoteip"`
}

type TokenValidateResponse struct {
	ErrorCodes []string `json:"error-codes"`
	Success    bool     `json:"success"`
	Action     string   `json:"action"`
	CData      int      `json:"cdata"`
}

func ValidateTurnstileToken(token string, secret string, remoteIp string) (bool, error) {
	req := TokenValidateRequest{
		Secret:   secret,
		Response: token,
		RemoteIP: remoteIp,
	}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	resp, err := http.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var res TokenValidateResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}
