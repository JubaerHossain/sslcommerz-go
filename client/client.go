package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	reader := strings.NewReader(`{
		"store_id": "impex60a637af47d85",
		"store_passwd": "impex60a637af47d85@ssl",
		"total_amount": "100.00",
		"currency":     "BDT",
		"tran_id":      "12345",
		"value_a":      "ref001_A",
		"value_b":      "ref002_B",
		"value_c":      "ref003_C",
		"cus_name":     "John Doe",
		"cus_add1":     "Dhaka, Bangladesh",
		"cus_city":     "Dhaka",
		"cus_postcode": "1000",
		"cus_country":  "Bangladesh",
		"cus_phone":    "01764824731",
		"cus_email":    "john.doe@example.com",
		"success_url":  "http://localhost:8080/success",
		"fail_url":     "http://localhost:8080/fail",
		"cancel_url":   "http://localhost:8080/cancel",
		"ipn_url":      "http://localhost:8080/ipn",
	}`)

	req, err := http.NewRequest(method, url, reader)
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
