package sslcommerz

import (
    "net/http"
    "encoding/json"
    "log"
)

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

func FailHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Payment Failed")
    http.Error(w, "Payment Failed", http.StatusBadRequest)
}

func CancelHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Payment Cancelled")
    w.Write([]byte("Payment Cancelled"))
}

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
