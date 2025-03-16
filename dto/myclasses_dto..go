package dto

// MyClasessResponse struct untuk response khusus menangani data batchgroup
type MyClasessResponse struct {
	ID        int  `json:"id"`
	TeacherID *int `json:"teacher_id"`
	BatchID   *int `json:"batch_id"`
	KursusID  *int `json:"kursus_id"`

	Teacher *responseUser   `json:"teacher"`
	Batch   *BatchResponse  `json:"batches"`
	Kursus  *responseKursus `json:"kursus"`
}
