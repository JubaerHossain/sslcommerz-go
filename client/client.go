package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/JubaerHossain/sslcommerz-go/config"
)

type Client struct {
	HttpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{},
	}
}

func (c *Client) MakeRequest(method, url string, payload interface{}) ([]byte, error) {
	var reqBody []byte
	var err error

	if payload != nil {
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(string(reqBody))

	fmt.Println("Making request to:", url)

	fmt.Println("")
	fmt.Println("")
	fmt.Println("method:", method)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")


	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.IS_SANDBOX == "true",
		},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
