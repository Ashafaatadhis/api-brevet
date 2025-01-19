package main

import (
	"new-brevet-be/config"
	"new-brevet-be/seed/seeder"
)

func main() {
	config.LoadEnv()
	// Inisialisasi database
	config.InitDB()

	// Seeder untuk data role dan user
	// seeder.RolesSeed()
	// seeder.HariSeed()
	// seeder.CategorySeed()
	seeder.PriceSeed()
	// seeder.KategoriGolonganSeed()
	// seeder.StatusPaymentSeed()
	// seeder.UserSeed()
	// seeder.KategoriSeed()
	// seeder.JenisKursusSeed()

}
