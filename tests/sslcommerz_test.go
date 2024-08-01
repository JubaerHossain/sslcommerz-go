// tests/sslcommerz_test.go
package tests

import (
	"fmt"
	"testing"
	"time"

	sslcommerzEntity "github.com/JubaerHossain/sslcommerz-go/pkg"
	"github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func TestInitiatePayment(t *testing.T) {
	paymentRequest := &sslcommerzEntity.PaymentRequest{
		TotalAmount:      103,
		Currency:         "BDT",
		TransactionID:    "SSLCZ_TEST_" + fmt.Sprintf("%d", time.Now().UnixNano()),
		SuccessURL:       "http://localhost:8080/success",
		FailURL:          "http://localhost:8080/fail",
		CancelURL:        "http://localhost:8080/cancel",
		IPNURL:           "http://localhost:8080/ipn",
		CustomerName:     "Test Customer",
		CustomerEmail:    "test@test.com",
		CustomerAddress1: "Dhaka",
		CustomerAddress2: "Dhaka",
		CustomerCity:     "Dhaka",
		CustomerState:    "Dhaka",
		CustomerPostcode: "1000",
		CustomerCountry:  "Bangladesh",
		CustomerPhone:    "01711111111",
		CustomerFax:      "01711111111",
		ShippingMethod:   "No",
		ShippingName:     "Store Test",
		ShippingAddress1: "Dhaka",
		ShippingAddress2: "Dhaka",
		ShippingCity:     "Dhaka",
		ShippingState:    "Dhaka",
		ShippingPostcode: "1000",
		ShippingCountry:  "Bangladesh",
		ValueA:           "ref001",
		ValueB:           "ref002",
		ValueC:           "ref003",
		ValueD:           "ref004",
		ProductName:      "Computer",
		ProductCategory:  "Goods",
		ProductProfile:   "physical-goods",
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
