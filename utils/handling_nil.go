package utils

// GetStringValue fungsi untuk menangani nil string
func GetStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// GetIntValue fungsi untuk menangani nil int
func GetIntValue(u *int) int {
	if u != nil {
		return *u
	}
	return 0
}
