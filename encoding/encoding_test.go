package encoding_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pragmaticengineering/netrunner/encoding"
)

type mockObject struct {
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		description string
		obj         mockObject
	}{
		{
			description: "encode a mock object",
			obj: mockObject{
				Name:  "John Doe",
				Email: "johndoe@example.com",
			},
		},
		// Add more test cases here if needed
	}

	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.description, func(t *testing.T) {
			// Create a mock response writer
			w := httptest.NewRecorder()

			// Define the test data
			status := http.StatusOK
			// Call the Encode function
			err := encoding.Encode(w, status, testcase.obj)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify the response writer
			resp := w.Result()
			defer resp.Body.Close()

			// Verify the status code
			if resp.StatusCode != status {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, status)
			}

			// Verify the content type header
			contentType := resp.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("unexpected content type: got %s, want application/json", contentType)
			}

			// Verify the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			ValidateObject := mockObject{}
			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&ValidateObject); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}
			if ValidateObject.Name != testcase.obj.Name || ValidateObject.Email != testcase.obj.Email {
				t.Errorf("unexpected response body: got %+v, want %+v", ValidateObject, testcase.obj)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		description string
		body        io.ReadCloser
		expected    mockObject
	}{
		{
			description: "decode a valid request body",
			body:        io.NopCloser(bytes.NewBufferString(`{"Name":"John Doe","Email":"johndoe@example.com"}`)),
			expected: mockObject{
				Name:  "John Doe",
				Email: "johndoe@example.com",
			},
		},
		// Add more test cases here if needed
	}

	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.description, func(t *testing.T) {
			// Call the Decode function
			result, err := encoding.Decode[mockObject](testcase.body)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify the decoded object
			if result.Name != testcase.expected.Name || result.Email != testcase.expected.Email {
				t.Errorf("unexpected decoded object: got %+v, want %+v", result, testcase.expected)
			}
		})
	}
}
