package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	sslcommerzEntity "github.com/JubaerHossain/sslcommerz-go/pkg"
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
	paymentRequest := &sslcommerzEntity.PaymentRequest{
		TotalAmount:      103,
		Currency:         "BDT",
		TransactionID:     "SSLCZ_TEST_" + GenerateUniqueID(),
		SuccessURL:        "http://localhost:8080/success",
		FailURL:           "http://localhost:8080/fail",
		CancelURL:         "http://localhost:8080/cancel",
		IPNURL:            "http://localhost:8080/ipn",
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
