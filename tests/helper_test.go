package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
)

// Register user buat keperluan testing
func registerUser(app *fiber.App, username, email, password string) {
	requestBody, _ := json.Marshal(validation.UserRegister{
		Name:     username,
		Username: username,
		Email:    email,
		Password: password,
		Nohp:     "08123456789",
		RoleID:   4,
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	log.Print(responseBody, "TEST INI")
}

// Login dan ambil token buat testing
func loginAndGetToken(app *fiber.App, username, email, password string) string {
	log.Print("mgew")
	registerUser(app, username, email, password)
	log.Print("mgew2")

	requestBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	log.Print(responseBody, "walawee")
	return responseBody["token"].(string)
}

func loginAsAdminAndGetToken(app *fiber.App) string {
	// 🔹 Register admin (kalau belum ada)
	registerBody, _ := json.Marshal(validation.UserRegister{
		Name:     "Admin",
		Username: "adminuser",
		Email:    "admin@example.com",
		Password: "adminpassword",
		Nohp:     "08123456789",
		RoleID:   1,
	})

	reqRegister := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	app.Test(reqRegister, -1) // Ignore response karena hanya register

	// 🔹 Login admin
	loginBody, _ := json.Marshal(map[string]string{
		"username": "adminuser",
		"password": "adminpassword",
	})
	reqLogin := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(reqLogin, -1)

	// 🔹 Ambil token dari response login
	var responseMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseMap)
	token := responseMap["token"].(string)

	log.Println("Admin Token:", token) // Debugging biar keliatan

	return token
}

func createBatchAndGetID(app *fiber.App) string {
	adminToken := loginAsAdminAndGetToken(app)
	log.Print("token nih : ", adminToken)
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

	resp, _ := app.Test(req, -1)
	// Parse response JSON
	var responseMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseMap)

	// Ambil data dari responseMap
	// Ambil data dari responseMap
	dataMap, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: Data tidak dalam format map[string]interface{}")
	}

	// Ambil id dari dataMap (float64)
	batchIDFloat, ok := dataMap["id"].(float64)
	if !ok {
		log.Fatal("Error: batchID bukan angka")
	}

	// Konversi float64 ke int
	batchIDInt := int(batchIDFloat)

	// Konversi int ke string
	batchIDStr := strconv.Itoa(batchIDInt)

	log.Println("Created Batch ID:", batchIDStr) // Logging ID batch untuk debugging

	return batchIDStr
}
func createCourseAndGetID(app *fiber.App) string {
	adminToken := loginAsAdminAndGetToken(app)

	// 🔹 Data batch yang akan dikirim
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

	}

	body, _ := json.Marshal(kursusRequest)

	// 🔹 Buat request POST /batch dengan token admin
	req := httptest.NewRequest(http.MethodPost, "/kursus", bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, _ := app.Test(req, -1)
	// Parse response JSON
	var responseMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseMap)

	// Ambil data dari responseMap
	dataMap, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: Data tidak dalam format map[string]interface{}")
	}

	// Ambil id dari dataMap (float64)
	kursusIDFloat, ok := dataMap["id"].(float64)
	if !ok {
		log.Fatal("Error: kursusID bukan angka")
	}

	// Konversi float64 ke int
	kursusIDInt := int(kursusIDFloat)

	// Konversi int ke string
	kursusIDStr := strconv.Itoa(kursusIDInt)

	log.Println("Created kursus ID:", kursusIDStr) // Logging ID batch untuk debugging

	return kursusIDStr
}

func createBatchMappingAndGetID(app *fiber.App) string {
	adminToken := loginAsAdminAndGetToken(app)

	// Ambil ID kursus yang berupa string
	kursusIDStr := createCourseAndGetID(app)
	kursusID, err := strconv.Atoi(kursusIDStr) // Konversi kursusID ke int
	if err != nil {
		log.Fatalf("Gagal mengonversi kursusID: %v", err)
	}

	// Ambil ID batch yang berupa string
	batchIDStr := createBatchAndGetID(app)
	batchID, err := strconv.Atoi(batchIDStr) // Konversi batchID ke int
	if err != nil {
		log.Fatalf("Gagal mengonversi batchID: %v", err)
	}
	// Membuat request POST untuk membuat batch mapping
	batchRequest := validation.CreateBatchMapping{
		KursusID: kursusID, // Gunakan kursusID yang sudah berupa int
		BatchID:  batchID,  // Gunakan batchID yang sudah berupa int
	}

	// Marshal request body
	body, _ := json.Marshal(batchRequest)
	req := httptest.NewRequest(http.MethodPost, "/batch-mapping", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, _ := app.Test(req, -1)
	// Parse response JSON
	var responseMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseMap)

	// Ambil data dari responseMap
	dataMap, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Error: Data tidak dalam format map[string]interface{}")
	}

	// Ambil id dari dataMap (float64)
	grBatchIDFloat, ok := dataMap["id"].(float64)
	if !ok {
		log.Fatal("Error: grBatchID bukan angka")
	}

	// Konversi float64 ke int
	grBatchIDInt := int(grBatchIDFloat)

	// Konversi int ke string
	grBatchIDStr := strconv.Itoa(grBatchIDInt)

	log.Println("Created kursus ID:", grBatchIDStr) // Logging ID batch untuk debugging

	return grBatchIDStr
}
