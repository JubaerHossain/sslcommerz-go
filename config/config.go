package config

import (
    "os"
)

var (
    StoreID     = os.Getenv("SSLCOMMERZ_STORE_ID")
    StorePass   = os.Getenv("SSLCOMMERZ_STORE_PASS")
    APIBaseURL  = "https://sandbox.sslcommerz.com"
)
