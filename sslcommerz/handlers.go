package sslcommerz

import (
    "encoding/json"
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
