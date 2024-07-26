package tests

import (
    "testing"
    "github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func TestInitiatePayment(t *testing.T) {
    paymentRequest := sslcommerz.PaymentRequest{
        TotalAmount: "100.00",
        Currency:    "BDT",
        TranID:      "12345",
        SuccessURL:  "http://example.com/success",
        FailURL:     "http://example.com/fail",
        CancelURL:   "http://example.com/cancel",
    }

    response, err := sslcommerz.InitiatePayment(paymentRequest)
    if err != nil {
        t.Fatalf("Error initiating payment: %v", err)
    }

    if response.Status != "SUCCESS" {
        t.Fatalf("Expected status to be SUCCESS, got %v", response.Status)
    }
}
