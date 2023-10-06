package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test data structure for data-driven tests
var testData = []struct {
	method       string
	path         string
	payload      interface{}
	expectedCode int
	expectedBody string
}{
	{
		method:       "POST",
		path:         "/?id=1",
		payload:      MyInfo{ID: "1", Name: "Alice", Age: 30},
		expectedCode: http.StatusCreated,
		expectedBody: "Created",
	},
	{
		method:       "GET",
		path:         "/?id=1",
		payload:      nil,
		expectedCode: http.StatusOK,
		expectedBody: "{\"number\":\"1\",\"name\":\"Alice\",\"age\":30}",
	},
	{
		method:       "PUT",
		path:         "/?id=1",
		payload:      MyInfo{ID: "1", Name: "Updated Alice", Age: 35},
		expectedCode: http.StatusOK,
		expectedBody: "OK",
	},
	{
		method:       "DELETE",
		path:         "/?id=1",
		payload:      nil,
		expectedCode: http.StatusOK,
		expectedBody: "OK",
	},
}

func TestMyServer_DataDriven(t *testing.T) {
	s := &MyServer{
		data: make(map[string]MyInfo),
	}

	for _, test := range testData {
		payload, _ := json.Marshal(test.payload)

		req, err := http.NewRequest(test.method, test.path, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Error creating %s request: %v", test.method, err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api-key", myAppSecret)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.ServeHTTP)
		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedCode {
			t.Errorf("Expected status %d for %s %s; got %d", test.expectedCode, test.method, test.path, rr.Code)
		}

		if !(string(test.expectedBody) == strings.Trim(rr.Body.String(), "\n")) {
			t.Errorf("Expected body '%v' for %s %s; got '%v'", test.expectedBody, test.method, test.path, rr.Body.String())
		}
	}
}
