# SSLCommerz-Go

A Go package for integrating with the SSLCommerz payment gateway. This package provides functions to initiate payments and handle responses, structured in a modular way for easy management and extension.

## Project Structure

```
sslcommerz-go/
├── config/
│   └── config.go
├── client/
│   └── client.go
├── sslcommerz/
│   └── sslcommerz.go
├── examples/
│   └── main.go
├── go.mod
└── go.sum
```

### Description of Files

- **config/config.go**: Handles configuration settings, such as storing the store ID and password.
- **client/client.go**: Manages HTTP requests to the SSLCommerz API.
- **sslcommerz/sslcommerz.go**: Implements payment processing functions, such as initiating payments.
- **examples/main.go**: Provides an example usage of the package, demonstrating how to initiate a payment request.

## Installation

1. **Clone the repository**:
   ```sh
   go get -u github.com/JubaerHossain/sslcommerz-go
   ```

2. **Set up Go modules**:
   ```sh
   go mod tidy
   ```

## Configuration

Set the following environment variables with your SSLCommerz credentials:

```sh
export SSLCOMMERZ_STORE_ID=your_store_id
export SSLCOMMERZ_STORE_PASS=your_store_password
```

## Usage

### Example Usage

The following example demonstrates how to initiate a payment request using the package.

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourusername/sslcommerz-go/sslcommerz"
)

func main() {
    paymentRequest := sslcommerz.PaymentRequest{
        TotalAmount: "100.00",
        Currency:    "BDT",
        TranID:      "12345",
        SuccessURL:  "http://example.com/success",
        FailURL:     "http://example.com/fail",
        CancelURL:   "http://example.com/cancel",
    }

    response, err := sslcommerz.InitiatePayment(paymentRequest)
    if err != nil {
        log.Fatalf("Error initiating payment: %v", err)
    }

    fmt.Printf("Payment Response: %+v\n", response)
}
```

### Run the Example

1. Set environment variables:
   ```sh
   export SSLCOMMERZ_STORE_ID=your_store_id
   export SSLCOMMERZ_STORE_PASS=your_store_password
   ```

2. Run the example:
   ```sh
   go run examples/main.go
   ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## Acknowledgements

Thanks to the SSLCommerz team for providing the API and documentation.

## Contact

For any questions or feedback, please open an issue or contact me at your-email@example.com.

---

This README provides a comprehensive guide for using the SSLCommerz-Go package, including setup, configuration, and example usage.