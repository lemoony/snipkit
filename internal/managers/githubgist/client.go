package githubgist

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"
)

var errAuth = errors.New("github unauthorized")

type rawResponse struct {
	hasUpdates    bool
	etag          string
	gistsResponse *[]rawGistsResponse
	rawContent    *[]byte
}

type rawGistsResponse struct {
	ID    string `json:"ID"`
	Files map[string]struct {
		Filename    string `json:"Filename"`
		ContentType string `json:"type"`
		Language    string `json:"language"`
		RawURL      string `json:"raw_url"`
	} `json:"files"`
	Public      bool   `json:"public"`
	Description string `json:"description"`
}

func (m Manager) checkToken(cfg GistConfig, token string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, cfg.apiURL(), nil)
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

	log.Trace().Msgf("Response status HEAD URL %s: %s", cfg.apiURL(), resp.Status)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.Errorf("unexpected status code from github: %s", resp.Status)
	}

	return true, nil
}

func (m Manager) getGists(cfg GistConfig, etag, token string) (rawResponse, error) {
	raw, err := m.getRawResponse(cfg.apiURL(), etag, token)
	if err != nil {
		return rawResponse{}, err
	}

	if !raw.hasUpdates {
		return rawResponse{hasUpdates: false}, nil
	}

	var response []rawGistsResponse
	err = json.Unmarshal(*raw.rawContent, &response)
	if err != nil {
		return rawResponse{}, err
	}

	return rawResponse{
		hasUpdates:    true,
		etag:          raw.etag,
		gistsResponse: &response,
	}, nil
}

func (m Manager) getRawGist(url, etag, token string) (rawResponse, error) {
	return m.getRawResponse(url, etag, token)
}

func (m Manager) getRawResponse(url, etag, token string) (rawResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return rawResponse{}, err
	}

	req.Header.Set("Accept", "application/vnd.github.VERSION.base64")

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}
	if etag != "" {
		req.Header.Set("If-None-Match", fmt.Sprintf(`"%s"`, etag))
	}

	resp, err := client.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Trace().Msgf("Response status %s URL %s: %s", req.Method, url, resp.Status)

	if err != nil {
		return rawResponse{}, err
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		if payload, err2 := ioutil.ReadAll(resp.Body); err != nil {
			panic(err2)
		} else {
			return rawResponse{}, errors.Wrap(errAuth, string(payload))
		}
	}

	if resp.StatusCode == http.StatusNotModified {
		return rawResponse{hasUpdates: false}, nil
	}

	if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return rawResponse{}, err
	} else {
		return rawResponse{
			hasUpdates: true,
			rawContent: &bytes,
			etag:       toStrongETag(resp.Header.Get("etag")),
		}, nil
	}
}

func toStrongETag(etag string) string {
	if strings.HasPrefix(etag, `W/"`) {
		return etag[3 : len(etag)-1]
	}
	return etag
}
