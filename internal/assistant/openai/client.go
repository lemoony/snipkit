package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/assistant/prompts"
	"github.com/lemoony/snipkit/internal/utils/httputil"
	"github.com/lemoony/snipkit/internal/utils/json"
)

type Client struct {
	config     Config
	httpClient httputil.HTTPClient
	history    []Message // Add history to store the conversation context
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
	log.Debug().Str("prompt", prompt).Msg("Starting query request")

	apiKey, err := c.apiKey()
	if err != nil {
		return "", err
	}

	reqBody, err := c.prepareRequest(prompt)
	if err != nil {
		return "", err
	}

	log.Trace().Str("request_body", string(reqBody)).Msg("Sending request to OpenAI")

	resp, err := c.sendRequest(apiKey, reqBody)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	return c.handleResponse(resp)
}

func (c *Client) prepareRequest(prompt string) ([]byte, error) {
	userMessage := Message{Role: "user", Content: prompt}

	if len(c.history) == 0 {
		systemMessage := Message{Role: "system", Content: prompts.DefaultPrompt}
		c.history = append(c.history, systemMessage)
	}

	c.history = append(c.history, userMessage)

	reqBody := Request{
		Model:    c.config.Model,
		Messages: c.history,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "Error marshaling request body")
	}

	return jsonBody, nil
}

func (c *Client) sendRequest(apiKey string, jsonBody []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("%s/v1/chat/completions", c.config.Endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error making request")
	}

	return resp, nil
}

func (c *Client) handleResponse(resp *http.Response) (string, error) {
	log.Debug().Int("statusCode", resp.StatusCode).Msg("Handling API response")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Error reading response body")
	}

	log.Trace().
		Int("status_code", resp.StatusCode).
		Str("response_body", string(body)).
		Msg("Received response from OpenAI")

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Error: received status code %d - Body: %s", resp.StatusCode, jsonutil.CompactJSON(body))
	}

	var openAIResp Response
	if err = json.Unmarshal(body, &openAIResp); err != nil {
		return "", errors.Wrap(err, "Error unmarshalling response body")
	}

	if len(openAIResp.Choices) > 0 {
		assistantMessage := openAIResp.Choices[0].Message
		c.history = append(c.history, assistantMessage)
		return assistantMessage.Content, nil
	}

	return "", errors.New("No response from OpenAI API.")
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
