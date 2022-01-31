package githubgist

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/phuslu/log"
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

func (m Manager) getGists(cfg GistConfig) ([]githubGistResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Trace().Msgf("Response status URL %s: %s", cfg.URL, resp.Status)

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
