package services

import (
	// Alias untuk crypto/rand
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand" // Tetap menggunakan nama default
	"new-brevet-be/utils"
	"os"
	"time"
)

// GenerateUniqueCode fungsi untuk generate unik code di harga (amount) nya
func GenerateUniqueCode(basePrice float64) float64 {
	// Menghasilkan angka acak antara 0 dan 999
	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)
	randomCode := randomGenerator.Intn(1000) // Angka acak antara 0 hingga 999

	// Gabungkan harga asli dengan kode acak
	finalPrice := basePrice + float64(randomCode)

	return finalPrice
}

// GenerateURLConfirm menghasilkan string acak yang unik untuk keperluan URL konfirmasi
func GenerateURLConfirm() (string, error) {
	// Menghasilkan byte acak
	randomBytes := make([]byte, 16) // 16 byte untuk keamanan yang baik

	_, err := cryptoRand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Mengubah byte acak menjadi string base64 yang aman untuk URL
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// Mengembalikan string acak sebagai kode konfirmasi
	return randomString, nil
}

// SendEmailConfirmAccount fungsi untuk mengatur pengirim email (kirim akun)
func SendEmailConfirmAccount(fullname string, emailuser string) error {

	// Pesan email dengan HTML
	subject := "Konfirmasi Pembayaran Tax Center Gunadarma"

	message := fmt.Sprintf(`
	<h3>Yth. Bapak/Ibu %s,</h3>
	<p>Terima kasih atas kepercayaan Bapak/Ibu telah melakukan pembelian pada website kami.</p>
	<p>Pembayaran Anda sudah kami verifikasi, silahkan login mengguankan akun yang sudah Bapak/Ibu daftar</p>
	 
	
	<p>Salam,<br>Tax Center Universitas Gunadarma</p>
	`, fullname)

	if err := utils.SendEmail(emailuser, subject, message); err != nil {
		return err
	}

	return nil
}

// SendEmailCodePayment fungsi untuk mengatur pengirim email (kirim kode pembayaran)
func SendEmailCodePayment(fullname, emailuser string, URLConfirm *string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	// Pesan email dengan HTML
	subject := "Konfirmasi Pembayaran Pelatihan Tax Center Universitas Gunadarma"

	message := fmt.Sprintf(`
	<h3>Yth. Bapak/Ibu %s,</h3>
	<p>Terima kasih atas kepercayaan Bapak/Ibu telah melakukan pembelian pada website kami.</p>

	<p>Untuk menyelesaikan proses pendaftaran, silakan melakukan pembayaran dengan menekan tombol dibawah</p>
	
    <a href="%s/confirmBayar/%s" style="display: inline-block; padding: 10px 20px; font-size: 16px; color: white; background-color: #007bff; text-decoration: none; border-radius: 5px;">Konfirmasi Pembayaran</a>
	 
	<p>Salam,<br>Tax Center Universitas Gunadarma</p>
`, fullname, frontendURL, *URLConfirm)

	if err := utils.SendEmail(emailuser, subject, message); err != nil {
		return err
	}

	return nil
}

// CreateUserAccount fungsi untuk membuat akun
// func CreateUserAccount(tx *gorm.DB, registration models.Registration) error {
// 	// Cek apakah email atau nomor HP sudah terdaftar
// 	var existingUser models.User
// 	if err := tx.Where("email = ? OR nohp = ?", registration.Email, registration.NoHP).First(&existingUser).Error; err == nil {
// 		// tetap lanjutkan ketika sudah ada gausah di create ulang
// 		return nil
// 	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
// 		return fmt.Errorf("failed to check existing user: %w", err)
// 	}

// 	// Hash password
// 	defaultPassword := os.Getenv("DEFAULT_PASSWORD_USER")
// 	hashedPassword, err := utils.HashPassword(defaultPassword)
// 	if err != nil {
// 		log.Println("Failed to hash password:", err)
// 		return fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Buat user baru
// 	userAccount := &models.User{
// 		Name:     registration.FullName,
// 		Username: registration.Email,
// 		Nohp:     registration.NoHP,
// 		RoleID:   4,
// 		Email:    registration.Email,
// 		Password: hashedPassword,
// 	}

// 	// Simpan user ke database
// 	if err := tx.Create(&userAccount).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("failed to register user: %w", err)
// 	}

// 	// Update kolom user_id di tabel registration
// 	if err := tx.Model(&registration).Update("user_id", userAccount.ID).Scan(&registration).Error; err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("failed to update registration: %w", err)
// 	}

// 	return nil // Jika tidak ada error
// }
