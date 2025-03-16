package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"new-brevet-be/validation"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostBatchMappingSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// Ambil ID kursus yang berupa string
	kursusIDStr := createCourseAndGetID(app)
	kursusID, err := strconv.Atoi(kursusIDStr) // Konversi kursusID ke int
	if err != nil {
		t.Fatalf("Gagal mengonversi kursusID: %v", err)
	}

	// Ambil ID batch yang berupa string
	batchIDStr := createBatchAndGetID(app)
	batchID, err := strconv.Atoi(batchIDStr) // Konversi batchID ke int
	if err != nil {
		t.Fatalf("Gagal mengonversi batchID: %v", err)
	}

	// Sekarang kursusID dan batchID sudah berupa int
	batchRequest := validation.CreateBatchMapping{
		KursusID: kursusID, // Gunakan kursusID yang sudah berupa int
		BatchID:  batchID,  // Gunakan batchID yang sudah berupa int
	}

	body, _ := json.Marshal(batchRequest)

	// 🔹 Buat request POST /batch dengan token admin
	req := httptest.NewRequest(http.MethodPost, "/batch-mapping", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	// 🔹 Cek apakah responsenya sesuai ekspektasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostBatchMappingFail(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// Ambil ID batch yang berupa string
	batchIDStr := createBatchAndGetID(app)
	batchID, err := strconv.Atoi(batchIDStr) // Konversi batchID ke int
	if err != nil {
		t.Fatalf("Gagal mengonversi batchID: %v", err)
	}

	// Sekarang kursusID dan batchID sudah berupa int
	batchRequest := validation.CreateBatchMapping{

		BatchID: batchID, // Gunakan batchID yang sudah berupa int
	}

	body, _ := json.Marshal(batchRequest)

	// 🔹 Buat request POST /batch dengan token admin
	req := httptest.NewRequest(http.MethodPost, "/batch-mapping", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := app.Test(req, -1)

	// 🔹 Cek apakah responsenya sesuai ekspektasi
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetBatchMapping(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// Menggunakan fungsi createBatchMappingAndGetID untuk membuat batch mapping
	createBatchMappingAndGetID(app)

	// 🔹 GET request untuk mengambil batch mapping
	getReq := httptest.NewRequest(http.MethodGet, "/batch-mapping", nil)
	getReq.Header.Set("Authorization", "Bearer "+adminToken)

	getResp, err := app.Test(getReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
}
func TestGetBatchMappingByID(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// 🔹 Login sebagai admin (role_id = 1) untuk mendapatkan token
	adminToken := loginAsAdminAndGetToken(app)

	// Menggunakan fungsi createBatchMappingAndGetID untuk membuat batch mapping
	ID := createBatchMappingAndGetID(app)

	// 🔹 GET request untuk mengambil batch mapping
	getReq := httptest.NewRequest(http.MethodGet, "/batch-mapping/"+ID, nil)
	getReq.Header.Set("Authorization", "Bearer "+adminToken)

	getResp, err := app.Test(getReq, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
}
