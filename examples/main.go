package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func main() {
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
        "cus_phone":    "017********",
        "cus_email":    "john.doe@example.com",
        "success_url":  "http://localhost:8080/success",
        "fail_url":     "http://localhost:8080/fail",
        "cancel_url":   "http://localhost:8080/cancel",
        "ipn_url":      "http://localhost:8080/ipn",
    }

    sslc := sslcommerz.NewSSLCommerz()
    response, err := sslc.InitiatePayment(paymentRequest)
    if err != nil {
        log.Fatalf("Error initiating payment: %v", err)
    }

    fmt.Printf("Payment Response: %+v\n", response)

    http.HandleFunc("/success", sslcommerz.SuccessHandler)
    http.HandleFunc("/fail", sslcommerz.FailHandler)
    http.HandleFunc("/cancel", sslcommerz.CancelHandler)
    http.HandleFunc("/ipn", sslcommerz.IPNHandler)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
