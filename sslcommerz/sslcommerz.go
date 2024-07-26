package sslcommerz

import (
    "encoding/json"
    "github.com/JubaerHossain/sslcommerz-go/client"
    "github.com/JubaerHossain/sslcommerz-go/config"
)

type PaymentRequest struct {
    StoreID     string `json:"store_id"`
    StorePass   string `json:"store_passwd"`
    TotalAmount string `json:"total_amount"`
    Currency    string `json:"currency"`
    TranID      string `json:"tran_id"`
    SuccessURL  string `json:"success_url"`
    FailURL     string `json:"fail_url"`
    CancelURL   string `json:"cancel_url"`
}

type PaymentResponse struct {
    Status string `json:"status"`
    // Add other response fields
}

func InitiatePayment(req PaymentRequest) (*PaymentResponse, error) {
    client := client.NewClient()
    req.StoreID = config.StoreID
    req.StorePass = config.StorePass

    url := config.APIBaseURL + "/gwprocess/v4/api.php"
    response, err := client.MakeRequest("POST", url, req)
    if err != nil {
        return nil, err
    }

    var paymentResponse PaymentResponse
    err = json.Unmarshal(response, &paymentResponse)
    if err != nil {
        return nil, err
    }

    return &paymentResponse, nil
}
