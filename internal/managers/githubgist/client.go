package githubgist

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"emperror.dev/errors"
	"github.com/phuslu/log"
)

var (
	errAuth = errors.New("github unauthorized")

	apiURLPattern = "https://api.%s/users/%s/gists"
)

type githubGistResponse struct {
	ID    string `json:"id"`
	Files map[string]struct {
		Filename    string `json:"filename"`
		ContentType string `json:"type"`
		Language    string `json:"language"`
		RawURL      string `json:"raw_url"`
	} `json:"files"`
}

func (m Manager) checkToken(cfg GistConfig, token string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, cfg.gistURL(), nil)
	if err != nil {
		return false, err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Trace().Msgf("Response status HEAD URL %s: %s", cfg.gistURL(), resp.Status)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("unexpected status code from github: %s", resp.Status)
	}

	return true, nil
}

func (m Manager) getGists(cfg GistConfig, token string) ([]githubGistResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, cfg.gistURL(), nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Trace().Msgf("Response status GET URL %s: %s", cfg.gistURL(), resp.Status)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		if payload, err2 := ioutil.ReadAll(resp.Body); err != nil {
			panic(err2)
		} else {
			return nil, errors.Wrap(errAuth, string(payload))
		}
	}

	var response []githubGistResponse

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (m Manager) getRawGist(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.VERSION.base64")

	resp, err := client.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("Response status URL %s: %s", url, resp.Status)

	if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}
