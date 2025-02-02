package dto

// PriceResponse struct untuk response khusus menangani data price
type PriceResponse struct {
	ID         int `json:"id"`
	Harga      int `json:"harga"`
	GolonganID int `json:"golongan_id"`

	KategoriGolongan KategoriGolonganResponse `json:"golongan"`
}
