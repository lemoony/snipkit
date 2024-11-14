package gemini

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
	restoreEnv := mockutil.MockAPIKeyEnv("GEMINI_API_KEY", "fake-api-key")
	defer restoreEnv()

	// Prepare a mock HTTP client.
	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gemini-1",
		Endpoint:  "https://api.google.com",
		APIKeyEnv: "GEMINI_API_KEY",
	}

	// Mock successful response from Google API.
	mockResponseBody := `{
		"candidates": [
			{
				"content": {
					"parts": [
						{"text": "Hello, world!"}
					]
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

	// Instantiate the client and inject the mock HTTP client.
	client := &Client{
		config:     config,
		httpClient: mockClient,
	}

	// Call the method and assert the result.
	result, err := client.Query("Hello?")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", result)

	// Assert that the HTTP client's Do method was called once.
	mockClient.AssertExpectations(t)
}

func TestClient_Query_APIKeyMissing(t *testing.T) {
	// Mock environment without API key.
	restoreEnv := mockutil.MockAPIKeyEnv("GEMINI_API_KEY", "")
	defer restoreEnv()

	config := Config{
		Model:     "gemini-1",
		Endpoint:  "https://api.google.com",
		APIKeyEnv: "GEMINI_API_KEY",
	}

	client := &Client{
		config: config,
	}

	// Call the method and assert the error for missing API key.
	_, err := client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "The environment variable GEMINI_API_KEY defined by apiKeyEnv is empty")
}

func TestClient_Query_HTTPError(t *testing.T) {
	// Mock API key environment variable.
	restoreEnv := mockutil.MockAPIKeyEnv("GEMINI_API_KEY", "fake-api-key")
	defer restoreEnv()

	// Prepare a mock HTTP client.
	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gemini-1",
		Endpoint:  "https://api.google.com",
		APIKeyEnv: "GEMINI_API_KEY",
	}

	mockClient.On("Do", mock.Anything).Return(&http.Response{}, errors.New("HTTP request failed"))

	client := &Client{
		config:     config,
		httpClient: mockClient,
	}

	// Call the method and assert the error from the HTTP client.
	_, err := client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP request failed")

	// Assert that the HTTP client's Do method was called.
	mockClient.AssertExpectations(t)
}

func TestClient_Query_InvalidResponse(t *testing.T) {
	// Mock API key environment variable.
	restoreEnv := mockutil.MockAPIKeyEnv("GEMINI_API_KEY", "fake-api-key")
	defer restoreEnv()

	mockClient := new(MockHTTPClient)

	config := Config{
		Model:     "gemini-1",
		Endpoint:  "https://api.google.com",
		APIKeyEnv: "GEMINI_API_KEY",
	}

	// Mock an invalid JSON response from the Google API.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
	}

	// Set up the mock to return the invalid response.
	mockClient.On("Do", mock.Anything).Return(mockResponse, nil)

	client := &Client{
		config:     config,
		httpClient: mockClient,
	}

	// Call the method and assert the error for invalid response.
	_, err := client.Query("Hello?")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error unmarshalling response body")

	// Assert that the HTTP client's Do method was called.
	mockClient.AssertExpectations(t)
}
