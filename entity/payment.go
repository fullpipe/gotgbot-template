package entity

import "gorm.io/gorm"

type Payment struct {
	gorm.Model

	UserID int64 `gorm:"index"`
	User   User

	Currency string
	Price    int64
	State    PaymentState

	TotalAmount             int64
	InvoicePayload          string
	ShippingOptionID        string
	TelegramPaymentChargeID string
	ProviderPaymentChargeID string
}

type PaymentState string

const (
	PendingPaymentState     PaymentState = "pending"
	PreCheckoutPaymentState PaymentState = "pre_checkout"
	SuccessPaymentState     PaymentState = "success"
	FailPaymentState        PaymentState = "fail"
)
