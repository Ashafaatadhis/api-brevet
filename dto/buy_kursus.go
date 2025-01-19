package dto

import (
	"new-brevet-be/models"
	"time"
)

// BuyKursusRequest adalah struct untuk request
type BuyKursusRequest struct {
	GroupBatchesID int `json:"group_batches_id" validate:"required,exists=group_batches.id,unique=group_batches.id"`
	JenisKursusID  int `json:"jenis_kursus_id" validate:"required,exists=jenis_kursus.id"`
}

// EditBuyKursus struct untuk response khusus menangani request
type EditBuyKursus struct {
	StatusPaymentID int `json:"status_payment_id" validate:"required,exists=status_payments.id"`
}

// BuykursusResponse adalah struct untuk response
type BuykursusResponse struct {
	ID              int       `json:"id"`
	GrBatchID       int       `json:"group_batches_id"`
	StatusPaymentID int       `json:"status_payment_id"`
	JenisKursusID   int       `json:"jenis_kursus_id"`
	UserID          *int      `json:"user_id"`
	URLConfirm      *string   `json:"url_confirm"`
	BuktiBayar      *string   `json:"bukti_bayar"`
	PriceID         int       `json:"price_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Price         price                `json:"price"`
	JenisKursus   models.JenisKursus   `json:"jenis_kursus"`
	GroupBatches  *GroupBatchResponse  `json:"group_batches"`
	User          *ResponseUser        `json:"user"`
	StatusPayment models.StatusPayment `json:"status_payment"`
}

type price struct {
	ID    int `json:"id"`
	Harga int `json:"harga"`
}
