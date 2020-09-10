# mock-http-api [![PkgGoDev](https://pkg.go.dev/badge/github.com/mkeeler/mock-http-api)](https://pkg.go.dev/github.com/mkeeler/mock-http-api)
Go helpers for mocking an HTTP API using stretchr/testify/mock

## Library Usage

```go
package mock_test

import (
   "encoding/json"
   "net/http"
   "testing"
   
   mockapi "github.com/mkeeler/mock-http-api"
)

// This test will pass as all the requisite API calls are made.
func TestMyAPI(t *testing.T) {
   m := mockapi.NewMockAPI(t)
   
   // This sets up an expectation that a GET request to /my/endpoint will be made and that it should
   // return a 200 status code with the provided map sent back JSON encoded as the body of the response
   call := m.WithJSONReply(mockapi.NewMockRequest("GET", "/my/endpoint"), 200, map[string]string{
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
```

## Code Generation

The code generator will create a new mock API type with helper methods for all the desired endpoints. These helpers
are meant to be more ergonomic to use that the raw `mock-http-api` module itself.

### Installing the code generator

```sh
go get github.com/mkeeler/mock-http-api/cmd/mock-api-gen
```

_Note that you may need to run this command with GO111MODULE=on if executing outside of your GOPATH_

### Using `mock-api-gen`

```sh
mock-api-gen -type MockMyAPI -endpoints ./endpoints.json -pkg myapi -output api.helpers.go
```

This command will take in the JSON file of endpoints and generate the desired type with helpers for mocking responses to each API.

The format of the endpoints file is:

```json
{
   "UpdateResource": {
      "Method": "POST",
      "Path": "/resource/%s",
      "PathParameters": ["resourceID"],
      "BodyType": "json",
      "ResponseType": "json",
      "Headers": true,
      "QueryParams": false
   }
}
```

Using this as input the following file would be generated:

```go
// Code generated by "mock-expect-gen -type MockAPI -pkg fakeapi -endpoints endpoints.json -output ./api.go"; DO NOT EDIT.

package fakeapi

import (
   "fmt"
   mockapi "github.com/mkeeler/mock-http-api"
)

type MockAPI struct {
   *mockapi.MockAPI
}

func NewMockAPI(t mockapi.TestingT) *MockAPI {
   return &MockAPI{
      MockAPI: mockapi.NewMockAPI(t),
   }
}

func (m *MockConsulAPI) UpdateResource(resourceID string, headers map[string]string, body map[string]interface{}, status int, reply interface{}) *mockapi.MockAPICall {
   req := mockapi.NewMockRequest("POST", fmt.Sprintf("/resource/%s", resourceID)).WithBody(body).WithHeaders(headers)

   return m.WithJSONReply(req, status, reply)
}
```

Then when you want to use this you would:

```
func TestFakeAPI(t *testing.T) {
   m := fakeapi.NewMockAPI(t)
   
   // Not necessary when the `t` passed into NewMockAPI supports a Cleanup method. (such as with the Go 1.14 testing.T type)
   defer m.Close()
   
   m.UpdateResource("some-id-here", 
      nil, 
      map[string]interface{"abc", "def"}, 
      200, 
      map[string]interface{"abc", "def", "added": true})
   
   httpServerURL := m.URL()
   
   // do something to cause the HTTP API call to happen here.
   
   // nothing else is necessary. Either the deferred m.Close or the automatic testing cleanup 
   // will assert that the required API calls were made.
}
```
 
#### Full Usage

```
Usage of mock-api-gen:
        mock-api-gen [flags] -type <type name> -endpoints <var name> [package]
Flags:
  -endpoints string
        File holding the endpoint configuration. (default "endpoints")
  -output string
        Output file name.
  -pkg string
        Name of the package to generate methods in
  -tag value
        Build tags the generated file should have. This may be specified multiple times.
  -type string
        Method receiver type the mock API helpers should be generated for
```