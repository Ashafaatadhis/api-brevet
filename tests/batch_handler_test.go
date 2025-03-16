package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"new-brevet-be/validation"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// **🟢 Test Buat Batch Berhasil**
func TestPostBatchSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// 🔹 Data batch yang akan dikirim
	batchRequest := validation.PostBatch{
		Judul:      "Batch Test",
		BukaBatch:  time.Now(),
		TutupBatch: time.Now().AddDate(0, 1, 0),
		JenisID:    1,  // Sesuaikan dengan ID yang valid
		Kuota:      10, // Minimal 0, pastikan ini diisi
		KelasID:    2,  // Sesuaikan dengan ID kelas yang valid
	}

	body, _ := json.Marshal(batchRequest)

	// 🔹 Buat request POST /batch dengan token admin
	req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	// 🔹 Cek apakah responsenya sesuai ekspektasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Buat Batch Gagal (Data Kosong)**
func TestPostBatchFail_EmptyData(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})
	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodPost, "/batch", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// **🟢 Test Ambil Detail Batch Berhasil**
func TestGetBatchDetailSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	batchID := createBatchAndGetID(app)

	req := httptest.NewRequest(http.MethodGet, "/batch/"+batchID, nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Ambil Detail Batch Gagal (ID Tidak Ada)**
func TestGetBatchDetailFail_NotFound(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/batch/9999", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// **🟢 Test Update Batch Berhasil**
func TestUpdateBatchSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	batchID := createBatchAndGetID(app)
	adminToken := loginAsAdminAndGetToken(app)
	updateRequest := validation.PostBatch{
		Judul:      "Updated Batch",
		BukaBatch:  time.Now(),
		TutupBatch: time.Now().AddDate(0, 1, 0),
		JenisID:    1,  // Sesuaikan dengan ID yang valid
		Kuota:      10, // Minimal 0, pastikan ini diisi
		KelasID:    2,  // Sesuaikan dengan ID kelas yang valid
	}
	body, _ := json.Marshal(updateRequest)

	req := httptest.NewRequest(http.MethodPut, "/batch/"+batchID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Update Batch Gagal (ID Tidak Valid)**
func TestUpdateBatchFail_InvalidID(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	updateRequest := validation.PostBatch{
		Judul:      "Updated Batch",
		BukaBatch:  time.Now(),
		TutupBatch: time.Now().AddDate(0, 1, 0),
		JenisID:    1,  // Sesuaikan dengan ID yang valid
		Kuota:      10, // Minimal 0, pastikan ini diisi
		KelasID:    2,  // Sesuaikan dengan ID kelas yang valid
	}
	body, _ := json.Marshal(updateRequest)

	req := httptest.NewRequest(http.MethodPut, "/batch/9999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// **🟢 Test Hapus Batch Berhasil**
func TestDeleteBatchSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	batchID := createBatchAndGetID(app)
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodDelete, "/batch/"+batchID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Hapus Batch Gagal (ID Tidak Ditemukan)**
func TestDeleteBatchFail_NotFound(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()
	adminToken := loginAsAdminAndGetToken(app)
	req := httptest.NewRequest(http.MethodDelete, "/batch/9999", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
