package sslcommerz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SuccessHandler handles the success response from SSLCommerz
func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	var payment map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Payment Success:", payment)
	w.Write([]byte("Payment Successful"))
}

// FailHandler handles the failure response from SSLCommerz
func FailHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment Failed")
	http.Error(w, "Payment Failed", http.StatusBadRequest)
}

// CancelHandler handles the cancellation response from SSLCommerz
func CancelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment Cancelled")
	w.Write([]byte("Payment Cancelled"))
}

// IPNHandler handles the IPN (Instant Payment Notification) response from SSLCommerz
func IPNHandler(w http.ResponseWriter, r *http.Request) {
	var payment map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("IPN Received:", payment)
	w.Write([]byte("IPN Received"))
}

func MakePaymentRequest(w http.ResponseWriter, r *http.Request) {
	paymentRequest := map[string]interface{}{
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
	}

	sslc := NewSSLCommerz()
	response, err := sslc.InitiatePayment(paymentRequest)
	if err != nil {
		log.Fatalf("Error initiating payment: %v", err)
		return 
	}

	returnURL := response["GatewayPageURL"].(string)
	http.Redirect(w, r, returnURL, http.StatusFound)
}
