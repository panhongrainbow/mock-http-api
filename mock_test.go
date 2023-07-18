package mockapi

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	// mockapi "github.com/mkeeler/mock-http-api"
)

// Test_String_Body is to check if there is an error when the body is not in JSON format but a simple string.
// It is performed to verify the fix.
func Test_String_Body(t *testing.T) {
	// Create a new instance of the MockAPI with the given testing.T instance
	m := NewMockAPI(t)

	// Set the filtered headers for the MockAPI
	m.SetFilteredHeaders([]string{
		"Accept-Encoding",
		"User-Agent",
		"Content-Length",
		"Content-Type",
	})

	// Create a new MockRequest instance with the HTTP method "POST" and path "/v1/ask"
	req := NewMockRequest(http.MethodPost, "/v1/ask").WithBody(
		[]byte("Golang is the most advanced language in the world."))

	// Create a new MockCall instance with a JSON reply containing the status code 200 and a map[string]string response
	call := m.WithJSONReply(req, 200, map[string]string{
		"response": "I agree with it.",
	})

	// Set the call to be executed only once
	call.Once()

	// Make an actual HTTP POST request to the MockAPI's URL with the specified path and request body
	resp, err := http.Post(fmt.Sprintf("%s/v1/ask", m.URL()), "text/plain",
		strings.NewReader("Golang is the most advanced language in the world."))
	require.NoError(t, err)

	// Defer the closing of the response body
	defer func() {
		_ = resp.Body.Close()
	}()

	// Create a new JSON decoder to decode the response body
	dec := json.NewDecoder(resp.Body)

	// Declare a variable to hold the decoded JSON output
	var output map[string]string

	// Decode the response body into the output variable
	err = dec.Decode(&output)
	require.NoError(t, err)

	// Assert that the "response" field in the output map is equal to "I agree with it."
	require.Equal(t, "I agree with it.", output["response"])
}

// This test will pass as all the requisite API calls are made.
func TestMyAPI(t *testing.T) {
	m := NewMockAPI(t)
	// http.Get will add both of the headers but we don't want to care about them.
	m.SetFilteredHeaders([]string{
		"Accept-Encoding",
		"User-Agent",
	})

	// This sets up an expectation that a GET request to /my/endpoint will be made and that it should
	// return a 200 status code with the provided map sent back JSON encoded as the body of the response
	call := m.WithJSONReply(NewMockRequest("GET", "/my/endpoint"), 200, map[string]string{
		"foo": "bar",
	})

	// This sets the call to be required to happen exactly once
	call.Once()

	// This makes the HTTP request to the mock HTTP server
	resp, err := http.Get(fmt.Sprintf("%s/my/endpoint", m.URL()))
	if err != nil {
		t.Fatalf("Error issuing GET of /my/endpoint: %v", err)
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	var output map[string]string
	if err := dec.Decode(&output); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if val, ok := output["foo"]; !ok || val != "bar" {
		t.Fatalf("Didn't get the expected response")
	}
}
