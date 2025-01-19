package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendEmail fungsi untuk mengatur pengirim email melalui smtp
func SendEmail(emailuser string, subject string, message string) error {

	// Konfigurasi SMTP
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	// Membuat pesan email
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser, "Tax Center Gunadarma")
	m.SetHeader("To", emailuser)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	// Mengirim email menggunakan SMTP
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = nil // Gunakan TLS

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
