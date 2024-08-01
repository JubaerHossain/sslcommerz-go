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

	fields := map[string]interface{}{
		"store_id":         payload.StoreID,
		"store_passwd":     payload.StorePass,
		"total_amount":     payload.TotalAmount,
		"currency":         payload.Currency,
		"tran_id":          payload.TransactionID,
		"success_url":      payload.SuccessURL,
		"fail_url":         payload.FailURL,
		"cancel_url":       payload.CancelURL,
		"ipn_url":          payload.IPNURL,
		"cus_name":         payload.CustomerName,
		"cus_email":        payload.CustomerEmail,
		"cus_add1":         payload.CustomerAddress1,
		"cus_add2":         payload.CustomerAddress2,
		"cus_city":         payload.CustomerCity,
		"cus_state":        payload.CustomerState,
		"cus_postcode":     payload.CustomerPostcode,
		"cus_country":      payload.CustomerCountry,
		"cus_phone":        payload.CustomerPhone,
		"cus_fax":          payload.CustomerFax,
		"shipping_method":  payload.ShippingMethod,
		"ship_name":        payload.ShippingName,
		"ship_add1":        payload.ShippingAddress1,
		"ship_add2":        payload.ShippingAddress2,
		"ship_city":        payload.ShippingCity,
		"ship_state":       payload.ShippingState,
		"ship_postcode":    payload.ShippingPostcode,
		"ship_country":     payload.ShippingCountry,
		"value_a":          payload.ValueA,
		"value_b":          payload.ValueB, // Now handled as a string
		"value_c":          payload.ValueC,
		"value_d":          payload.ValueD,
		"product_name":     payload.ProductName,
		"product_category": payload.ProductCategory,
		"product_profile":  payload.ProductProfile,
	}

	for key, value := range fields {
		if err := writer.WriteField(key, fmt.Sprint(value)); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %v", key, err)
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
