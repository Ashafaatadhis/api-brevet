package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"new-brevet-be/validation"

	"github.com/stretchr/testify/assert"
)

// **🟢 Test Register Berhasil**
func TestAuthRegisterSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	requestBody, _ := json.Marshal(validation.UserRegister{
		Name:     "Test User2",
		Username: "testuser2s",
		Email:    "test2@example.com",
		Password: "password123",
		Nohp:     "081234456789",
		RoleID:   4,
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode response body
	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	assert.Equal(t, "User registered successfully", responseBody["message"])
}

// **❌ Test Register Gagal - Password Kosong**
func TestAuthRegisterFail_EmptyPassword(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	requestBody, _ := json.Marshal(map[string]string{
		"name":     "Test User",
		"username": "testuser",
		"email":    "test@example.com",
		"password": "",
		"nohp":     "08123456789",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// **❌ Test Register Gagal - Email Sudah Terdaftar**
func TestAuthRegisterFail_EmailAlreadyExists(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// Register user pertama
	requestBody, _ := json.Marshal(validation.UserRegister{
		Name:     "User 1",
		Username: "user1",
		Email:    "duplicate@example.com",
		Password: "password123",
		Nohp:     "08123456789",
		RoleID:   4,
	})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	// Register user kedua dengan email yang sama
	requestBody2, _ := json.Marshal(validation.UserRegister{
		Name:     "User 2",
		Username: "user2",
		Email:    "duplicate@example.com", // Email yang sama
		Password: "password456",
		Nohp:     "08123456780",
		RoleID:   4,
	})
	req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody2))
	req2.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req2)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// **🟢 Test Login Berhasil**
func TestAuthLoginSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// Register dulu biar user ada
	registerUser(app, "testuser", "test@example.com", "password123")

	// Login
	requestBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Cek apakah response punya token
	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NotEmpty(t, responseBody["token"])
}

// **❌ Test Login Gagal - Password Salah**
func TestAuthLoginFail_WrongPassword(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	// Register dulu biar user ada
	registerUser(app, "testuser", "test@example.com", "password123")

	// Login dengan password salah
	requestBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// **🟢 Test Ambil Data User (`/me`)**
func TestAuthMeSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})
	app := setupApp()

	// Register & Login buat dapet token
	token := loginAndGetToken(app, "testuser", "test@example.com", "password123")
	log.Print("CEK")

	log.Print(token, "CEK TOKEN")

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// **❌ Test Ambil Data User Gagal - Token Tidak Valid**
func TestAuthMeFail_InvalidToken(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// **🟢 Test Logout Berhasil**
func TestAuthLogoutSuccess(t *testing.T) {
	t.Cleanup(func() {
		cleanupDatabase(testDB) // 🔥 Bersihkan data setelah tes ini selesai
	})

	app := setupApp()

	token := loginAndGetToken(app, "testuser", "test@example.com", "password123")
	log.Print(token, "TEST TOOKEN")
	req := httptest.NewRequest(http.MethodDelete, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
