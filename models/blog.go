package models

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Blog model untuk tabel blogs
type Blog struct {
	ID        int    `gorm:"primaryKey;type:int unsigned;"`
	Slug      string `gorm:"type:varchar(255);uniqueIndex;not null;"`
	Judul     string `gorm:"type:varchar(255);not null;"`
	Deskripsi string `gorm:"type:text;not null;"`
	Content   string `gorm:"type:longtext;not null;"`
	Gambar    string `gorm:"type:varchar(255);"`
}

// BeforeCreate hook untuk generate slug unik sebelum insert ke database
func (b *Blog) BeforeCreate(tx *gorm.DB) (err error) {
	b.Slug, err = generateUniqueSlug(tx, b.Judul)
	return
}

// generateUniqueSlug buat slug unik berdasarkan judul
func generateUniqueSlug(tx *gorm.DB, judul string) (string, error) {
	// Konversi judul jadi slug dasar
	slug := strings.ToLower(strings.ReplaceAll(judul, " ", "-"))
	originalSlug := slug
	var count int64

	// Cek apakah slug sudah ada di database
	err := tx.Model(&Blog{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return "", err
	}

	fmt.Printf("Slug pertama: %s, Count: %d\n", slug, count) // Debugging

	// Jika slug sudah ada, tambahkan angka di belakangnya
	i := 1
	for count > 0 {
		slug = fmt.Sprintf("%s-%d", originalSlug, i)
		err := tx.Model(&Blog{}).Where("slug = ?", slug).Count(&count).Error
		if err != nil {
			return "", err
		}

		fmt.Printf("Mencoba slug: %s, Count: %d\n", slug, count) // Debugging
		i++
	}

	fmt.Printf("Slug final: %s\n", slug) // Debugging
	return slug, nil
}
