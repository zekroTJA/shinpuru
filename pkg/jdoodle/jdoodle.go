// Package jdoodle provides an API wrapper for
// the jdoodle execute and credit-spent REST API.
package jdoodle

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	apiRoot = "https://api.jdoodle.com/v1/"
)

// Wrapper provides a request client holding
// the authentication credentials.
type Wrapper struct {
	clientId     string
	clientSecret string
}

// NewWrapper creates a new instance of wrapper
// using the given clientId and clientSecret.
func NewWrapper(clientId, clientSecret string) *Wrapper {
	return &Wrapper{clientId, clientSecret}
}

// ExecuteScript executes the given script of the
// given lang using the jdoodle execute API endpoint.
func (jd *Wrapper) ExecuteScript(lang, script string) (res *ExecResponse, err error) {
	payload := &execRequestBody{
		credentialsBody: &credentialsBody{
			ClientID:     jd.clientId,
			ClientSecret: jd.clientSecret,
		},
		Language: lang,
		Script:   script,
	}

	res = new(ExecResponse)
	err = request("execute", payload, res)

	return
}

// CreditsSpent returns the number of spent API
// credits today.
func (jd *Wrapper) CreditsSpent() (res *CreditsResponse, err error) {
	payload := &credentialsBody{
		ClientID:     jd.clientId,
		ClientSecret: jd.clientSecret,
	}

	res = new(CreditsResponse)
	err = request("credit-spent", payload, res)

	return
}

func request(endpoint string, body interface{}, res interface{}) (err error) {
	buf := bytes.NewBuffer([]byte{})
	err = json.NewEncoder(buf).Encode(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", apiRoot+endpoint, buf)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if httpRes.StatusCode >= 400 {
		if httpRes.ContentLength == 0 {
			err = errors.New(httpRes.Status)
			return
		}

		errBody := new(responseError)
		err = json.NewDecoder(httpRes.Body).Decode(errBody)
		if err != nil {
			err = errors.New(httpRes.Status)
		} else {
			err = errors.New(errBody.Error)
		}

		return
	}

	err = json.NewDecoder(httpRes.Body).Decode(res)

	return
}
