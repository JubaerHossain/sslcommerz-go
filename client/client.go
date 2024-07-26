package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/JubaerHossain/sslcommerz-go/config"
)

// Client represents an HTTP client
type Client struct {
	HttpClient *http.Client
}

// NewClient creates a new instance of Client
func NewClient() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IS_SANDBOX == "true", // Adjust based on environment
		},
	}
	return &Client{
		HttpClient: &http.Client{Transport: tr},
	}
}

// MakeRequest sends an HTTP request with the given method, URL, and payload
func (c *Client) MakeRequest(method, url string, payload map[string]interface{}) ([]byte, error) {
	// var reqBody []byte
	var err error

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for key, value := range payload {
		if err := writer.WriteField(key, fmt.Sprint(value)); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %v", key, err)
		}
	}

	writer.Close()

	// Send the POST request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp2, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp2.Body.Close()

	body, err2 := io.ReadAll(resp2.Body)
	if err2 != nil {
		return nil, fmt.Errorf("error reading response body: %v", err2)
	}
	return body, nil
}
