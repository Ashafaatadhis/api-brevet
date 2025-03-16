package dto

// CreateBlogRequest digunakan saat membuat blog baru
type CreateBlogRequest struct {
	Judul     string `json:"judul" validate:"required"`
	Deskripsi string `json:"deskripsi" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

// UpdateBlogRequest digunakan saat mengupdate blog (semua field optional)
type UpdateBlogRequest struct {
	Judul     *string `json:"judul,omitempty" validate:"omitempty,max=255"`
	Deskripsi *string `json:"deskripsi,omitempty" validate:"omitempty"`
	Content   *string `json:"content,omitempty" validate:"omitempty"`
}

// BlogResponse digunakan untuk mengembalikan data blog
type BlogResponse struct {
	ID        int     `json:"id"`
	Slug      string  `json:"slug"`
	Judul     string  `json:"judul"`
	Deskripsi string  `json:"deskripsi"`
	Content   string  `json:"content"`
	Gambar    *string `json:"gambar,omitempty"` // <-- Hanya dihapus dari JSON jika nil
}
