package validation

import (
	"new-brevet-be/dto"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// InitValidator menginisialisasi validator dengan kustom validation
func InitValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		validate.RegisterValidation("exists", ValidateExists)
		validate.RegisterValidation("unique", ValidateUnique)
		validate.RegisterValidation("unique_except", ValidateUniqueExcept(""))
		validate.RegisterStructValidation(CreateBuyKursusUniqueValidation, dto.BuyKursusRequest{})

	})
	return validate
}
