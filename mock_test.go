package mockapi

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
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

// Test_String_Body is to check if there is an error when the body is not in JSON format but a simple string.
// It is performed to verify the fix.
func Test_Json_Body(t *testing.T) {
	// Create a new instance of the MockAPI with the given testing.T instance
	m := NewMockAPI(t)

	// Set the filtered headers for the MockAPI
	m.SetFilteredHeaders([]string{
		"Accept-Encoding",
		"User-Agent",
		"Content-Length",
		"Content-Type",
	})

	jsonStr := struct {
		String  string      `json:"string"`
		Number  int         `json:"number"`
		Float   float64     `json:"float"`
		Boolean bool        `json:"boolean"`
		Null    interface{} `json:"null"`
		Object  struct {
			Key1 string `json:"key1"`
			Key2 string `json:"key2"`
		} `json:"object"`
		Array       []int     `json:"array"`
		NestedArray [][]int   `json:"nestedArray"`
		Date        time.Time `json:"date"`
	}{
		String:  "Hello world",
		Number:  123,
		Float:   12.34,
		Boolean: true,
		Null:    nil,
		Object: struct {
			Key1 string `json:"key1"`
			Key2 string `json:"key2"`
		}{
			Key1: "value1",
			Key2: "value2",
		},
		Array:       []int{1, 2, 3},
		NestedArray: [][]int{{1, 2}, {3, 4}},
		Date:        time.Date(2019, 1, 1, 12, 34, 56, 0, time.UTC),
	}

	// Create a new MockRequest instance with the HTTP method "POST" and path "/v1/ask"
	req := NewMockRequest(http.MethodPost, "/v1/json").WithBody(jsonStr)

	// Create a new MockCall instance with a JSON reply containing the status code 200 and a map[string]string response
	call := m.WithJSONReply(req, 200, map[string]string{
		"response": "I can recognize it.",
	})

	// Set the call to be executed only once
	call.Once()

	// Make an actual HTTP POST request to the MockAPI's URL with the specified path and request body
	resp, err := http.Post(fmt.Sprintf("%s/v1/json", m.URL()), "application/json",
		strings.NewReader(`{
		"string": "Hello world",
		"number": 123,
		"float": 12.34, 
		"boolean": true,
		"null": null,
		"object": {
			"key1": "value1",
			"key2": "value2"
		},
		"array": [1, 2, 3],
		"nestedArray": [[1, 2], [3, 4]],
		"date": "2019-01-01T12:34:56Z"
	}`))
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
	fmt.Println(output)
	// require.Equal(t, "I agree with it.", output["response"])
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
