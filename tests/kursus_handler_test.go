package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"new-brevet-be/validation"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// **🟢 Test Buat Kursus Berhasil**
func TestPostKursusSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// 🔹 Data kursus yang akan dikirim
	kursusRequest := validation.PostKursus{
		Judul:            "Kursus Go Programming",
		DeskripsiSingkat: "Belajar Go dari dasar hingga mahir",
		Deskripsi:        "Materi lengkap mengenai Go, dimulai dari sintaks dasar hingga implementasi di dunia nyata.",
		Pembelajaran:     "Kursus ini akan mengajarkan tentang cara menggunakan Go dengan berbagai contoh praktis.",
		Diperoleh:        "Setelah menyelesaikan kursus, peserta akan menguasai Go.",
		CategoryID:       1,                                                        // ID kategori yang valid
		ThumbnailKursus:  "http://example.com/thumbnail.jpg",                       // URL thumbnail
		StartDate:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),   // Tanggal mulai kursus
		EndDate:          time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC), // Tanggal selesai kursus
		StartTime:        "09:00:00",                                               // Jam mulai
		EndTime:          "12:00:00",                                               // Jam selesai
		HariID:           []uint{1, 2},                                             // ID hari yang relevan, misalnya Senin dan Selasa
		// Opsional, bisa disesuaikan dengan data yang ada di `models.Hari`
	}

	body, _ := json.Marshal(kursusRequest)

	// 🔹 Buat request POST /kursus dengan token admin
	req := httptest.NewRequest(http.MethodPost, "/kursus", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)
	log.Print(resp, "INI RESPONN")
	// 🔹 Cek apakah responsenya sesuai ekspektasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Buat Kursus Gagal (Data Kosong)**
func TestPostKursusFail_EmptyData(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodPost, "/kursus", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// **🟢 Test Ambil Detail Kursus Berhasil**
func TestGetKursusDetailSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	kursusID := createCourseAndGetID(app)

	req := httptest.NewRequest(http.MethodGet, "/kursus/"+kursusID, nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Ambil Detail Kursus Gagal (ID Tidak Ada)**
func TestGetKursusDetailFail_NotFound(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/kursus/9999", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// **🟢 Test Update Kursus Berhasil**
func TestUpdateKursusSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	kursusID := createCourseAndGetID(app)
	adminToken := loginAsAdminAndGetToken(app)
	updateRequest := validation.PostKursus{
		Judul:            "Updated Kursus Go Programming",
		DeskripsiSingkat: "Update materi kursus Go",
		Deskripsi:        "Materi tambahan tentang Go, dengan fokus pada fitur baru.",
		Pembelajaran:     "Mempelajari Go versi terbaru.",
		Diperoleh:        "Peserta dapat menguasai fitur terbaru dari Go.",
		CategoryID:       1,                                                        // ID kategori yang valid
		ThumbnailKursus:  "http://example.com/updated-thumbnail.jpg",               // URL thumbnail
		StartDate:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),   // Tanggal mulai kursus
		EndDate:          time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC), // Tanggal selesai kursus
		StartTime:        "09:00:00",                                               // Jam mulai
		EndTime:          "12:00:00",                                               // Jam selesai
		HariID:           []uint{1, 2},                                             // ID hari yang relevan
		Hari:             nil,                                                      // Opsional, bisa disesuaikan dengan data yang ada di `models.Hari`
	}
	body, _ := json.Marshal(updateRequest)

	req := httptest.NewRequest(http.MethodPut, "/kursus/"+kursusID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Update Kursus Gagal (ID Tidak Valid)**
func TestUpdateKursusFail_InvalidID(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	updateRequest := validation.PostKursus{
		Judul:            "Updated Kursus Go Programming",
		DeskripsiSingkat: "Update materi kursus Go",
		Deskripsi:        "Materi tambahan tentang Go, dengan fokus pada fitur baru.",
		Pembelajaran:     "Mempelajari Go versi terbaru.",
		Diperoleh:        "Peserta dapat menguasai fitur terbaru dari Go.",
		CategoryID:       1,                                                        // ID kategori yang valid
		ThumbnailKursus:  "http://example.com/updated-thumbnail.jpg",               // URL thumbnail
		StartDate:        time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),   // Tanggal mulai kursus
		EndDate:          time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC), // Tanggal selesai kursus
		StartTime:        "09:00:00",                                               // Jam mulai
		EndTime:          "12:00:00",                                               // Jam selesai
		HariID:           []uint{1, 2},                                             // ID hari yang relevan
		Hari:             nil,                                                      // Opsional, bisa disesuaikan dengan data yang ada di `models.Hari`
	}
	body, _ := json.Marshal(updateRequest)

	req := httptest.NewRequest(http.MethodPut, "/kursus/9999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// **🟢 Test Hapus Kursus Berhasil**
func TestDeleteKursusSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	kursusID := createCourseAndGetID(app)
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodDelete, "/kursus/"+kursusID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Hapus Kursus Gagal (ID Tidak Ditemukan)**
func TestDeleteKursusFail_NotFound(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodDelete, "/kursus/9999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
