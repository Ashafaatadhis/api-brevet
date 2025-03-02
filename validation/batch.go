package validation

import (
	"time"
)

// PostBatch struct untuk response khusus menangani data batch
type PostBatch struct {
	Judul      string    `validate:"required,min=3,max=100" json:"judul"`            // Judul wajib diisi, minimal 3 karakter, maksimal 100
	BukaBatch  time.Time `validate:"required" json:"buka_batch"`                     // Tanggal buka batch wajib diisi
	TutupBatch time.Time `validate:"required,gtefield=BukaBatch" json:"tutup_batch"` // Tanggal tutup batch wajib >= BukaBatch
	JenisID    int       `json:"jenis_id" validate:"required"`
	Kuota      int       `json:"kuota" validate:"required,gte=0"`
	KelasID    int       `json:"kelas_id" validate:"required"`
}
