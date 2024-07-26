package sslcommerz

import (
    "encoding/json"
    "errors"
    "github.com/JubaerHossain/sslcommerz-go/client"
    "github.com/JubaerHossain/sslcommerz-go/config"
    "net/url"
    "os"
)

type SSLCommerz struct {
    storeID         string
    storePass       string
    sslcSubmitURL   string
    sslcValidationURL string
    sslcMode        string
    sslcData        map[string]interface{}
}

func NewSSLCommerz() *SSLCommerz {
    sslc := &SSLCommerz{
        storeID:   config.StoreID,
        storePass: config.StorePass,
    }

    if os.Getenv("APP_ENV") == "production" {
        sslc.sslcMode = "securepay"
        sslc.sslcSubmitURL = "https://securepay.sslcommerz.com/gwprocess/v4/api.php"
        sslc.sslcValidationURL = "https://securepay.sslcommerz.com/validator/api/validationserverAPI.php"
    } else {
        sslc.sslcMode = "sandbox"
        sslc.sslcSubmitURL = "https://sandbox.sslcommerz.com/gwprocess/v4/api.php"
        sslc.sslcValidationURL = "https://sandbox.sslcommerz.com/validator/api/validationserverAPI.php"
    }

    return sslc
}

func (s *SSLCommerz) InitiatePayment(postData map[string]interface{}) (map[string]interface{}, error) {
    postData["store_id"] = s.storeID
    postData["store_passwd"] = s.storePass

    client := client.NewClient()
    response, err := client.MakeRequest("POST", s.sslcSubmitURL, postData)
    if err != nil {
        return nil, err
    }

    var result map[string]interface{}
    err = json.Unmarshal(response, &result)
    if err != nil {
        return nil, err
    }

    if status, ok := result["status"].(string); ok && status == "SUCCESS" {
        return result, nil
    }

    return nil, errors.New("failed to initiate payment")
}

func (s *SSLCommerz) ValidateTransaction(tranID, amount, currency string, postData map[string]interface{}) (bool, error) {
    if tranID == "" || amount == "" || currency == "" {
        return false, errors.New("invalid transaction data")
    }

    postData["store_id"] = s.storeID
    postData["store_passwd"] = s.storePass

    validationURL := s.sslcValidationURL + "?val_id=" + url.QueryEscape(postData["val_id"].(string)) +
        "&store_id=" + url.QueryEscape(s.storeID) +
        "&store_passwd=" + url.QueryEscape(s.storePass) +
        "&v=1&format=json"

    client := client.NewClient()
    response, err := client.MakeRequest("GET", validationURL, nil)
    if err != nil {
        return false, err
    }

    var result map[string]interface{}
    err = json.Unmarshal(response, &result)
    if err != nil {
        return false, err
    }

    if status, ok := result["status"].(string); ok && (status == "VALID" || status == "VALIDATED") {
        return true, nil
    }

    return false, errors.New("transaction validation failed")
}

// Additional methods for refund, transaction check, etc. can be implemented similarly
