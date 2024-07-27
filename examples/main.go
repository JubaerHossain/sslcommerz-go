package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func main() {

	http.HandleFunc("GET /payment", MakePaymentRequest)
	http.HandleFunc("POST /success", SuccessHandler)
	http.HandleFunc("POST /fail", FailHandler)
	http.HandleFunc("POST /cancel", CancelHandler)
	http.HandleFunc("POST /ipn", IPNHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// SuccessHandler handles the success response from SSLCommerz
func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Example: accessing specific form data
	paymentID := r.FormValue("tran_id")
	amount := r.FormValue("amount")
	status := r.FormValue("status")

	log.Println("Payment Success:")
	log.Println("Transaction ID:", paymentID)
	log.Println("Amount:", amount)
	log.Println("Status:", status)

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

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	tranID := r.FormValue("tran_id")
	if tranID == "" {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	data := make(map[string]interface{})
	for k, v := range r.Form {
		data[k] = v[0]
	}
	sslc := sslcommerz.NewSSLCommerz()
	amount := "100"
	// valid, err := sslc.ValidateTransaction(tranID, amount, "BDT", data)

	valid, err := sslc.ValidateTransaction(tranID, fmt.Sprintf("%.2f", amount), "BDT", data)
	if err != nil || !valid {
		// Save the updated master order status to the database
		log.Println("Validation failed")
		w.Write([]byte("Validation Fail"))
		return
	}

	fmt.Println("IPN Received")
	// Save the updated master order status to the database

	w.Write([]byte("IPN Received"))
}

func GenerateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// MakePaymentRequest handles the payment request to SSLCommerz.
func MakePaymentRequest(w http.ResponseWriter, r *http.Request) {
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

	sslc := sslcommerz.NewSSLCommerz()
	response, err := sslc.InitiatePayment(paymentRequest)
	if err != nil {
		log.Printf("Error initiating payment: %v", err)
		http.Error(w, "Internal Server Error: Payment initiation failed", http.StatusInternalServerError)
		return
	}

	status, ok := response["status"].(string)
	if !ok || status != "SUCCESS" {
		log.Printf("Payment initiation failed with status: %v", status)
		http.Error(w, "Payment initiation failed", http.StatusInternalServerError)
		return
	}

	gatewayURL, ok := response["GatewayPageURL"].(string)
	if !ok || gatewayURL == "" {
		log.Printf("Invalid or missing GatewayPageURL in response: %v", response)
		http.Error(w, "Invalid GatewayPageURL in response", http.StatusInternalServerError)
		return
	}
	// Ensure redirect occurs after sending JSON response
	http.Redirect(w, r, gatewayURL, http.StatusSeeOther)
}
