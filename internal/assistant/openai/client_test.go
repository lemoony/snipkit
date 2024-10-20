package openai

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestClient_Query_Success(t *testing.T) {
	// Mock API key environment variable.
	restoreEnv := mockutil.MockAPIKeyEnv("OPENAI_API_KEY", "fake-api-key")
	defer restoreEnv()

	// Prepare a mock HTTP httpClient.
	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gpt-3.5-turbo",
		Endpoint:  "https://api.openai.com",
		APIKeyEnv: "OPENAI_API_KEY",
	}

	// Mock successful response from OpenAI API.
	mockResponseBody := `{
		"choices": [
			{
				"message": {
					"role": "assistant",
					"content": "Hello, world!"
				}
			}
		]
	}`

	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(mockResponseBody)),
	}

	// Set up the mock to return the mock response when called.
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	// Instantiate the httpClient and inject the mock HTTP httpClient.
	client, err := NewClient(WithConfig(config), WithHTTPClient(mockClient))
	assert.NoError(t, err)

	// Call the method and assert the result.
	result, err := client.Query("Hello?")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", result)

	// Assert that the HTTP httpClient's Do method was called once.
	mockClient.AssertExpectations(t)
}

func TestClient_Query_APIKeyMissing(t *testing.T) {
	// Mock environment without API key.
	restoreEnv := mockutil.MockAPIKeyEnv("OPENAI_API_KEY", "")
	defer restoreEnv()

	config := Config{
		Model:     "gpt-3.5-turbo",
		Endpoint:  "https://api.openai.com",
		APIKeyEnv: "OPENAI_API_KEY",
	}

	client, err := NewClient(WithConfig(config))
	assert.NoError(t, err)

	// Call the method and assert the error for missing API key.
	_, err = client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The environment variable OPENAI_API_KEY defined by apiKeyEnv is empty")
}

func TestClient_Query_HTTPError(t *testing.T) {
	// Mock API key environment variable.
	restoreEnv := mockutil.MockAPIKeyEnv("OPENAI_API_KEY", "fake-api-key")
	defer restoreEnv()

	// Prepare a mock HTTP httpClient.
	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gpt-3.5-turbo",
		Endpoint:  "https://api.openai.com",
		APIKeyEnv: "OPENAI_API_KEY",
	}

	// Mock an error response from the HTTP httpClient.
	mockClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("HTTP request failed"))

	client, err := NewClient(WithConfig(config), WithHTTPClient(mockClient))
	assert.NoError(t, err)

	// Call the method and assert the error from the HTTP httpClient.
	_, err = client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP request failed")

	// Assert that the HTTP httpClient's Do method was called.
	mockClient.AssertExpectations(t)
}

func TestClient_Query_InvalidResponse(t *testing.T) {
	// Mock API key environment variable.
	restoreEnv := mockutil.MockAPIKeyEnv("OPENAI_API_KEY", "fake-api-key")
	defer restoreEnv()

	// Prepare a mock HTTP httpClient.
	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gpt-3.5-turbo",
		Endpoint:  "https://api.openai.com",
		APIKeyEnv: "OPENAI_API_KEY",
	}

	// Mock an invalid JSON response from the OpenAI API.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
	}

	// Set up the mock to return the invalid response.
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	client, err := NewClient(WithConfig(config), WithHTTPClient(mockClient))
	assert.NoError(t, err)

	// Call the method and assert the error for invalid response.
	_, err = client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response body")

	// Assert that the HTTP httpClient's Do method was called.
	mockClient.AssertExpectations(t)
}
