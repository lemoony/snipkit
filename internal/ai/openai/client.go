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

	"github.com/lemoony/snipkit/internal/ai/prompts"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Client struct {
	system *system.System
	config Config
	cache  cache.Cache
}

func NewClient(options ...Option) (*Client, error) {
	manager := &Client{}
	for _, o := range options {
		o.apply(manager)
	}
	return manager, nil
}

func (c *Client) Query(prompt string) string {
	apiKey := os.Getenv("SNIPKIT_OPENAI_API_KEY")
	if apiKey == "" {
		panic("Please set the OPENAI_API_KEY environment variable.")
	}

	systemMessage := Message{Role: "system", Content: prompts.DefaultPrompt}
	userMessage := Message{Role: "user", Content: prompt}

	reqBody := Request{
		Model:    c.config.Model,
		Messages: []Message{systemMessage, userMessage},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		panic(errors.Wrap(err, "Error marshaling request body"))
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("%s/v1/chat/completions", c.config.Endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(errors.Wrap(err, "Error creating request"))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		panic(errors.Wrap(err, "Error making request"))
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "Error reading response body"))
	}

	if resp.StatusCode != http.StatusOK {
		panic(errors.Errorf("Error: received status code %d - Body: %s", resp.StatusCode, string(body)))
	}

	var openAIResp Response
	if err = json.Unmarshal(body, &openAIResp); err != nil {
		panic(errors.Wrap(err, "Error unmarshalling response body"))
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content
	}
	panic(errors.New("No response from OpenAI API."))
}
