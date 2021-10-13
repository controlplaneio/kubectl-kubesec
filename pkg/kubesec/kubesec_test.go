package kubesec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewHTTPTestServer - Returns a new instances of httptest server with pre-determined response
func NewHTTPTestServer(code int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		_, err := w.Write([]byte(body))
		if err != nil {
			fmt.Println("error writing HTTP response")
		}
	}))
}

func TestKubesecClient_ScanDefinition(t *testing.T) {
	testCases := []struct {
		name             string
		mockResponseCode int
		mockResponseBody string
		wantErr          bool
		timeOutSec       int
		wantScore        int
		errorString      string
	}{
		{"Test non 200 response", 500, `validResp`, true, 5, 6, "got 500 response"},
		{"Test 200 response but with empty result", 200, ``, true, 5, 6, "failed to scan definition"},
		{"Test 200 response but with invalid result json", 200, `"foo": "bar"`, true, 5, 6, "invalid character"},
		{"Test timeout", 200, `"foo": "bar"`, true, 0, 6, "context deadline exceeded"},
		{"Test valid response with valid score", 200, `[{"score": 6,"scoring": {}}]`, false, 10, 6, ""},
		{"Test with an Error field set to non-empty string", 200, `[{"score": 60,"error": "dummyError"}]`, false, 10, 60, ""},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var b bytes.Buffer
			server := NewHTTPTestServer(tc.mockResponseCode, tc.mockResponseBody)

			defer server.Close()

			kc := NewClient(server.URL, tc.timeOutSec)

			// our http handler does not care about the data being sent
			_, err := b.WriteString(`dummyPost`)

			if err != nil {
				t.Fatal(err)
			}

			result, err := kc.ScanDefinition(b)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorString)
			} else {
				require.NoError(t, err)
				assert.Equal(t, result[0].Score, tc.wantScore)
			}
		})
	}
}

func TestKubeSecResults(t *testing.T) {

	responseWithOutError := `
	[{
	"score": -36,
	"scoring": {
		"critical": [{
				"id": "Privileged",
				"selector": "containers[] .securityContext .privileged == true",
				"reason": "Privileged containers can allow almost completely unrestricted host access",
				"points": -30
			},
			{
				"id": "HostNetwork",
				"selector": ".spec .hostNetwork == true",
				"reason": "Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter",
				"points": -9
			}
		]
	}
	}]
	`
	validOutput := "kubesec.io score: -36\n-----------------\nCritical\n1. containers[] .securityContext .privileged == true\n" +
		"Privileged containers can allow almost completely unrestricted host access\n2. .spec .hostNetwork == true\n" +
		"Sharing the host's network namespace permits processes in the pod to communicate with " +
		"processes bound to the host's loopback adapter\n-----------------\n"

	testCases := []struct {
		name     string
		response string
		want     string
		wantErr  bool
	}{
		{"response with error field not set", responseWithOutError, validOutput, false},
		{"response with error field set", `[{"score": -36,"error":"dummyError"}]`, "", true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var results KubeSecResults
			var b bytes.Buffer

			err := json.Unmarshal([]byte(tc.response), &results)
			if err != nil {
				t.Fatal(err)
			}
			err = results.Dump(&b)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, b.String())
			}
		})
	}
}
