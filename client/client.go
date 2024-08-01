package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"github.com/JubaerHossain/sslcommerz-go/config"
	sslcommerzEntity "github.com/JubaerHossain/sslcommerz-go/pkg"
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
func (c *Client) MakeRequest(method, url string, payload *sslcommerzEntity.PaymentRequest) ([]byte, error) {
	// var reqBody []byte
	var err error

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Tag.Get("json")

		if err := writer.WriteField(fieldName, fmt.Sprint(field.Interface())); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %v", fieldName, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}
	// Send the POST request
	client := &http.Client{
		Timeout: 60 * time.Second,
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
