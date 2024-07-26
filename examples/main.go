package main

import (
	"log"
	"net/http"

	"github.com/JubaerHossain/sslcommerz-go/sslcommerz"
)

func main() {

	http.HandleFunc("/payment", sslcommerz.MakePaymentRequest)
	http.HandleFunc("/success", sslcommerz.SuccessHandler)
	http.HandleFunc("/fail", sslcommerz.FailHandler)
	http.HandleFunc("/cancel", sslcommerz.CancelHandler)
	http.HandleFunc("/ipn", sslcommerz.IPNHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
