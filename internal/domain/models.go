package domain

import "github.com/gofrs/uuid"

type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ProductCode string  `json:"productCode"`
}

type PaymentMethod struct {
	DownPaymentAmount float64 `json:"downPaymentAmount"`
	Installments      int     `json:"installments"`
	TotalAmount       float64 ``
}

type PurchaseSummary struct {
	ProductInfo Product       `json:"productInfo"`
	PaymentInfo PaymentMethod `json:"paymentInfo"`
	ID          uuid.UUID     `json:"id"`
}

type PurchaseSummaryResponse struct {
	TotalAmount  float64 `json:"totalAmount"`
	InterestRate float64 `json:"InterestRate"`
	Installments int     `json:"installments"`
}
