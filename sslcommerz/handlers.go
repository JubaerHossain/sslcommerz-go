package sslcommerz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

func GenerateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func MakePaymentRequest(w http.ResponseWriter, r *http.Request) {
	// Create the payment request payload
	paymentRequest := map[string]interface{}{
		"total_amount":        "103",
		"currency":            "BDT",
		"tran_id":             "SSLCZ_TEST_" + GenerateUniqueID(),
		"success_url":         "http://localhost:8080/success",
		"fail_url":            "http://localhost:8080/fail",
		"cancel_url":          "http://localhost:8080/cancel",
		"ipn_url":             "http://localhost:8080/ipn",
		"emi_option":          "1",
		"emi_max_inst_option": "9",
		"emi_selected_inst":   "9",
		"cus_name":            "Test Customer",
		"cus_email":           "test@test.com",
		"cus_add1":            "Dhaka",
		"cus_add2":            "Dhaka",
		"cus_city":            "Dhaka",
		"cus_state":           "Dhaka",
		"cus_postcode":        "1000",
		"cus_country":         "Bangladesh",
		"cus_phone":           "01711111111",
		"cus_fax":             "01711111111",
		"shipping_method":     "No",
		"ship_name":           "Store Test",
		"ship_add1":           "Dhaka",
		"ship_add2":           "Dhaka",
		"ship_city":           "Dhaka",
		"ship_state":          "Dhaka",
		"ship_postcode":       "1000",
		"ship_country":        "Bangladesh",
		"value_a":             "ref001",
		"value_b":             "ref002",
		"value_c":             "ref003",
		"value_d":             "ref004",
		"product_name":        "Computer",
		"product_category":    "Goods",
		"product_profile":     "physical-goods",
	}

	// Initialize SSLCommerz client
	sslc := NewSSLCommerz()
	response, err := sslc.InitiatePayment(paymentRequest)
	if err != nil {
		log.Printf("Error initiating payment: %v", err)
		http.Error(w, "Internal Server Error: Payment initiation failed", http.StatusInternalServerError)
		return
	}

	// Check if the payment initiation was successful
	status, ok := response["status"].(string)
	if !ok || status != "SUCCESS" {
		log.Printf("Payment initiation failed with status: %v", status)
		http.Error(w, "Payment initiation failed", http.StatusInternalServerError)
		return
	}

	// Extract the gateway URL from the response
	gatewayURL, ok := response["GatewayPageURL"].(string)
	if !ok || gatewayURL == "" {
		log.Printf("Invalid or missing GatewayPageURL in response: %v", response)
		http.Error(w, "Invalid GatewayPageURL in response", http.StatusInternalServerError)
		return
	}

	// Send success response and redirect to the gateway URL
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Payment initiated successfully",
		"gateway_url": gatewayURL,
	})
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error: Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	// Redirect to the gateway URL
	http.Redirect(w, r, gatewayURL, http.StatusSeeOther)
}
