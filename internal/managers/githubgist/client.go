package githubgist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"
)

var (
	errAuth       = errors.New("github unauthorized")
	errUnexpected = errors.New("unexpected status code from github")
)

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

func (m Manager) checkToken(cfg GistConfig, token string) bool {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, cfg.apiURL(), nil)
	if err != nil {
		panic(err)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Trace().Msgf("Response status HEAD URL %s: %s", cfg.apiURL(), resp.Status)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.Wrap(errUnexpected, resp.Status))
	}

	return true
}

func (m Manager) getGists(cfg GistConfig, etag, token string) rawResponse {
	raw := m.getRawResponse(cfg.apiURL(), etag, token)
	if !raw.hasUpdates {
		return rawResponse{hasUpdates: false}
	}

	var response []rawGistsResponse
	err := json.Unmarshal(*raw.rawContent, &response)
	if err != nil {
		panic(err)
	}

	return rawResponse{
		hasUpdates:    true,
		etag:          raw.etag,
		gistsResponse: &response,
	}
}

func (m Manager) getRawGist(url, etag, token string) rawResponse {
	return m.getRawResponse(url, etag, token)
}

func (m Manager) getRawResponse(url, etag, token string) rawResponse {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		panic(err)
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

	if err != nil {
		panic(err)
	}

	log.Trace().Msgf("Response status %s URL %s: %s", req.Method, url, resp.Status)

	if resp.StatusCode == http.StatusNotModified {
		return rawResponse{hasUpdates: false}
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		if payload, err2 := io.ReadAll(resp.Body); err2 != nil {
			panic(err2)
		} else {
			panic(errors.Wrap(errAuth, string(payload)))
		}
	}

	if bytes, err2 := io.ReadAll(resp.Body); err2 != nil {
		panic(err2)
	} else {
		return rawResponse{
			hasUpdates: true,
			rawContent: &bytes,
			etag:       toStrongETag(resp.Header.Get("etag")),
		}
	}
}

func toStrongETag(etag string) string {
	if strings.HasPrefix(etag, `W/"`) {
		return etag[3 : len(etag)-1]
	}
	return etag
}
