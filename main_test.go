// server_test.go

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func sendRequestAndGetResponse(s *MyServer, method, path, apiKey string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("api-key", apiKey) // Set the API key
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)
	return recorder
}

func TestGetMyInfo(t *testing.T) {
	initialData := map[string]MyInfo{
		"1": {ID: "1", Name: "Alice", Age: 30},
		"2": {ID: "2", Name: "Bob", Age: 35},
	}
	s := &MyServer{
		data: initialData,
	}

	tests := []struct {
		id         string
		statusCode int
		apiKey     string
	}{
		{"1", http.StatusOK, "myAppSecret12254"},
		{"3", http.StatusUnauthorized, "invalidApiKey"},
	}

	for _, test := range tests {
		path := "/?id=" + test.id
		recorder := sendRequestAndGetResponse(s, http.MethodGet, path, test.apiKey)

		if recorder.Code != test.statusCode {
			t.Errorf("for ID %s, expected status %v, got %v", test.id, test.statusCode, recorder.Code)
		}

		if test.statusCode == http.StatusOK {
			expectedBody, _ := json.Marshal(initialData[test.id])
			if string(expectedBody) != strings.Trim(recorder.Body.String(), "\n") {
				t.Errorf("for ID %s, expected response %v, got %v", test.id, string(expectedBody), recorder.Body.String())
			}
		}
	}
}
