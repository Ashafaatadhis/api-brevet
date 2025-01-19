package validation

// CreateBatchMapping struct untuk response khusus menangani data batch-mapping
type CreateBatchMapping struct {
	KursusID uint `json:"kursus_id" validate:"required"`
	BatchID  uint `json:"batch_id" validate:"required"`
}
