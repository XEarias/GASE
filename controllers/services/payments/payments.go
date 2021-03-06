package payments

//Payment define a basic payment struct
type Payment struct {
	ID        string
	OrderID   int
	Method    string //webcheckout or CC data or tokenized
	Price     float32
	Currency  string
	Tax       float32
	ExtraData map[string]interface{}
	Token     string
	Country   string
}

//CreditCardPayment defines a struct to Credit Cards Charges
type CreditCardPayment struct {
	Payment
	CreditCard struct {
		CCV     string
		ExpDate string
		Number  string
	}
}

//WebCheckoutPayment defines a struct to WebChekcouts Payments
type WebCheckoutPayment struct {
	Payment
	Status  string
	PayerID string
}
