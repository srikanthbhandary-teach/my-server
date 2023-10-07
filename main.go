/*
Package client provides a client for interacting with the MyInfo server.

This package is maintained by Srikanth Bhandary.
*/
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client represents the client for the MyInfo server.
type Client struct {
	baseURL string
	apiKey  string
}

// NewClient creates a new Client with the specified baseURL and apiKey.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// CreateMyInfo sends a POST request to create a new MyInfo entity.
func (c *Client) CreateMyInfo(id, name string, age int) error {
	url := fmt.Sprintf("%s/?id=%s", c.baseURL, id)

	info := map[string]interface{}{
		"number": id,
		"name":   name,
		"age":    age,
	}

	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("error marshaling request payload: %v", err)
	}

	return c.sendRequest("POST", url, payload)
}

// GetMyInfo sends a GET request to retrieve a MyInfo entity by its ID.
func (c *Client) GetMyInfo(id string) (*MyInfo, error) {
	url := fmt.Sprintf("%s/?id=%s", c.baseURL, id)

	responseBody, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var info MyInfo
	if err := json.Unmarshal(responseBody, &info); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &info, nil
}

// UpdateMyInfo sends a PUT request to update an existing MyInfo entity.
func (c *Client) UpdateMyInfo(id, name string, age int) error {
	url := fmt.Sprintf("%s/?id=%s", c.baseURL, id)

	info := map[string]interface{}{
		"number": id,
		"name":   name,
		"age":    age,
	}

	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("error marshaling request payload: %v", err)
	}

	return c.sendRequest("PUT", url, payload)
}

// DeleteMyInfo sends a DELETE request to delete a MyInfo entity by its ID.
func (c *Client) DeleteMyInfo(id string) error {
	url := fmt.Sprintf("%s/?id=%s", c.baseURL, id)

	return c.sendRequest("DELETE", url, nil)
}

func (c *Client) sendRequest(method, url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating %s request: %v", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending %s request: %v", method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return responseBody, nil
}
