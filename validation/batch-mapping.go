package validation

// CreateBatchMapping struct untuk response khusus menangani data batch-mapping
type CreateBatchMapping struct {
	KursusID int `json:"kursus_id" validate:"required"`
	BatchID  int `json:"batch_id" validate:"required"`
}
