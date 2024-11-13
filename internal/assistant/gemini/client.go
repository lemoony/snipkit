package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/assistant/prompts"
	"github.com/lemoony/snipkit/internal/utils/httputil"
)

type Client struct {
	config       Config
	httpClient   httputil.HTTPClient
	contentParts []ContentParts
}

func NewClient(options ...Option) (*Client, error) {
	manager := &Client{
		httpClient: &http.Client{},
	}
	for _, o := range options {
		o.apply(manager)
	}
	return manager, nil
}

func (c *Client) Query(prompt string) (string, error) {
	apiKey, err := c.apiKey()
	if err != nil {
		return "", err
	}

	reqBody, err := c.prepareRequest(prompt)
	if err != nil {
		return "", err
	}

	resp, err := c.sendRequest(apiKey, reqBody)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	return c.handleResponse(resp)
}

func (c *Client) prepareRequest(prompt string) ([]byte, error) {
	c.contentParts = append(c.contentParts, ContentParts{
		Role:  "user",
		Parts: []TextPart{{Text: prompt}},
	})

	reqBody := Request{
		SystemInstruction: Instruction{Parts: TextPart{Text: prompts.DefaultPrompt}},
		Contents:          c.contentParts,
		SafetySettings: []SafetySetting{
			{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "BLOCK_NONE"},
			{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "BLOCK_NONE"},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling request body")
	}

	return jsonBody, nil
}

func (c *Client) sendRequest(apiKey string, jsonBody []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		"POST",
		fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.config.Endpoint, c.config.Model, apiKey),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error making request")
	}

	return resp, nil
}

func (c *Client) handleResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Error: received status code %d - Body: %s", resp.StatusCode, string(body))
	}

	var googleAIResp Response
	if err = json.Unmarshal(body, &googleAIResp); err != nil {
		return "", errors.Wrap(err, "Error unmarshalling response body")
	}

	if len(googleAIResp.Candidates) > 0 && len(googleAIResp.Candidates[0].Content.Parts) > 0 {
		c.contentParts = append(c.contentParts, googleAIResp.Candidates[0].Content)
		return googleAIResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", errors.New("No response from Google API.")
}

func (c *Client) apiKey() (string, error) {
	if apiKeyEnv := c.config.APIKeyEnv; len(apiKeyEnv) > 0 {
		if apiKey := os.Getenv(apiKeyEnv); len(apiKey) > 0 {
			return apiKey, nil
		}
		return "", errors.Errorf("The environment variable %s defined by apiKeyEnv is empty", apiKeyEnv)
	}

	return "", errors.Errorf("No environment variable specified for property 'apiKeyEnv' in config")
}