// tests/sslcommerz_test.go
package tests

import (
    "testing"
    "github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func TestInitiatePayment(t *testing.T) {
    paymentRequest := map[string]interface{}{
        "total_amount": "100.00",
        "currency":     "BDT",
        "tran_id":      "12345",
        "success_url":  "http://localhost:8080/success",
        "fail_url":     "http://localhost:8080/fail",
        "cancel_url":   "http://localhost:8080/cancel",
        "ipn_url":      "http://localhost:8080/ipn",
        "cus_name":     "John Doe",
        "cus_add1":     "123 Main St",
        "cus_city":     "Dhaka",
        "cus_postcode": "1000",
        "cus_country":  "Bangladesh",
        "cus_phone":    "017xxxxxxxx",
        "cus_email":    "john.doe@example.com",
    }

    sslc := sslcommerz.NewSSLCommerz()
    response, err := sslc.InitiatePayment(paymentRequest)
    if err != nil {
        t.Fatalf("Error initiating payment: %v", err)
    }

    if status, ok := response["status"].(string); !ok || status != "SUCCESS" {
        t.Fatalf("Expected status to be SUCCESS, got %v", status)
    }
}
