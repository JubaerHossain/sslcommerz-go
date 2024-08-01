package sslcommerzEntity

type PaymentRequest struct {
	StoreID          string  `json:"store_id"`
	StorePass        string  `json:"store_passwd"`
	TotalAmount      float64 `json:"total_amount"`
	Currency         string  `json:"currency"`
	TransactionID    string  `json:"tran_id"`
	SuccessURL       string  `json:"success_url"`
	FailURL          string  `json:"fail_url"`
	CancelURL        string  `json:"cancel_url"`
	IPNURL           string  `json:"ipn_url"`
	CustomerName     string  `json:"cus_name"`
	CustomerEmail    string  `json:"cus_email"`
	CustomerAddress1 string  `json:"cus_add1"`
	CustomerAddress2 string  `json:"cus_add2"`
	CustomerCity     string  `json:"cus_city"`
	CustomerState    string  `json:"cus_state"`
	CustomerPostcode string  `json:"cus_postcode"`
	CustomerCountry  string  `json:"cus_country"`
	CustomerPhone    string  `json:"cus_phone"`
	CustomerFax      string  `json:"cus_fax"`
	ShippingMethod   string  `json:"shipping_method"`
	ShippingName     string  `json:"ship_name"`
	ShippingAddress1 string  `json:"ship_add1"`
	ShippingAddress2 string  `json:"ship_add2"`
	ShippingCity     string  `json:"ship_city"`
	ShippingState    string  `json:"ship_state"`
	ShippingPostcode string  `json:"ship_postcode"`
	ShippingCountry  string  `json:"ship_country"`
	ValueA           string  `json:"value_a"`
	ValueB           string  `json:"value_b"`
	ValueC           string  `json:"value_c"`
	ValueD           string  `json:"value_d"`
	ProductName      string  `json:"product_name"`
	ProductCategory  string  `json:"product_category"`
	ProductProfile   string  `json:"product_profile"`
}
