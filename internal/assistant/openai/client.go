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

	"github.com/lemoony/snipkit/internal/assistant/prompts"
	"github.com/lemoony/snipkit/internal/utils/httputil"
	jsonutil "github.com/lemoony/snipkit/internal/utils/json"
)

type Client struct {
	config     Config
	httpClient httputil.HTTPClient
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

	systemMessage := Message{Role: "system", Content: prompts.DefaultPrompt}
	userMessage := Message{Role: "user", Content: prompt}

	reqBody := Request{
		Model:    c.config.Model,
		Messages: []Message{systemMessage, userMessage},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", errors.Wrap(err, "Error marshaling request body")
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("%s/v1/chat/completions", c.config.Endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", errors.Wrap(err, "Error creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Error making request")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Error reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Error: received status code %d - Body: %s", resp.StatusCode, jsonutil.CompactJSON(body))
	}

	var openAIResp Response
	if err = json.Unmarshal(body, &openAIResp); err != nil {
		return "", errors.Wrap(err, "Error unmarshalling response body")
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content, nil
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
