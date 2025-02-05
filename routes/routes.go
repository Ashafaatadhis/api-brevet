package routes

import (
	"new-brevet-be/dto"
	"new-brevet-be/handlers"
	"new-brevet-be/middlewares"
	"new-brevet-be/policy"
	"new-brevet-be/validation"

	"github.com/gofiber/fiber/v2"
)

// Setup untuk handling routing
func Setup(v1 fiber.Router) {
	// auth
	auth := v1.Group("/auth")
	auth.Get("/me", middlewares.AuthMiddleware(), handlers.Me)
	auth.Post("/register", validation.Validate[validation.UserRegister](), middlewares.UserUniqueCheck[validation.UserRegister], handlers.Register())
	auth.Post("/login", validation.Validate[validation.UserLogin](), handlers.Login())
	auth.Delete("/logout", middlewares.AuthMiddleware(), handlers.Logout())

	// user-setting
	v1.Put("/user-setting", middlewares.AuthMiddleware(), validation.Validate[validation.UserSetting](),
		middlewares.UserUniqueCheck[validation.UserSetting], handlers.UpdateUserProfile())
	v1.Delete("/user-setting/avatar", middlewares.AuthMiddleware(),
		handlers.DeleteAvatar)

	// change password
	v1.Put("/change-password", middlewares.AuthMiddleware(), validation.Validate[validation.ChangePassword](), handlers.ChangePassword)

	// Manage user (khusus admin)
	v1.Get("/manage-user", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.GetManageUser)
	v1.Get("/manage-user/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.GetDetailManageUser)
	v1.Post("/manage-user", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostManageUser](), middlewares.UserUniqueCheck[validation.PostManageUser], handlers.PostManageUser)
	v1.Put("/manage-user/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.UpdateManageUser](), middlewares.UserUniqueCheck[validation.UpdateManageUser],
		handlers.UpdateManageUser)
	v1.Delete("/manage-user/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteManageUser)

	// Kursus
	v1.Get("/kursus", handlers.GetKursus)
	v1.Get("/kursus/:id", handlers.GetDetailKursus)
	v1.Post("/kursus", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostKursus](), middlewares.KursusUniqueCheck[validation.PostKursus],
		handlers.PostKursus)
	v1.Put("/kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostKursus](), middlewares.KursusUniqueCheck[validation.PostKursus],
		handlers.UpdateKursus)
	v1.Delete("/kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteKursus)

	// Batch
	v1.Get("/batch", handlers.GetBatch)
	v1.Get("/batch/:id", handlers.GetDetailBatch)
	v1.Post("/batch", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostBatch](), middlewares.BatchUniqueCheck[validation.PostBatch],
		handlers.PostBatch)
	v1.Put("/batch/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.PostBatch](), middlewares.BatchUniqueCheck[validation.PostBatch],
		handlers.UpdateBatch)
	v1.Delete("/batch/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteBatch)

	// batch-mapping
	v1.Get("/batch-mapping", handlers.GetAllBatchMappping)
	v1.Get("/batch-mapping/:id", handlers.GetDetailBatchMappping)
	v1.Post("/batch-mapping", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.CreateBatchMapping](),
		handlers.CreateBatchMapping)
	v1.Put("/batch-mapping/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[validation.CreateBatchMapping](),
		handlers.EditBatchMapping)
	v1.Delete("/batch-mapping/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		handlers.DeleteBatchMapping)

	//  mapping pengajar (khusus admin)
	v1.Put("/mapping-pengajar/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"admin"}),
		validation.Validate[dto.MappingPengajarRequest](), middlewares.GuruUniqueCheck[dto.MappingPengajarRequest],
		handlers.MappingPengajar)

	// manage guru (khusus helpdesk)
	v1.Get("/manage-guru", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetManageGuru)
	v1.Get("/manage-guru/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetDetailManageGuru)
	v1.Post("/manage-guru", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		validation.Validate[validation.PostManageGuru](), middlewares.UserUniqueCheck[validation.PostManageGuru], handlers.PostManageGuru)
	v1.Put("/manage-guru/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		validation.Validate[validation.UpdateManageGuru](), middlewares.UserUniqueCheck[validation.UpdateManageGuru], handlers.UpdateManageGuru)
	v1.Delete("/manage-guru/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.DeleteManageGuru)

	// registration (user registration to brevet)
	v1.Get("/registration", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetAllRegistration)
	v1.Get("/registration/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetDetailRegistration)
	v1.Post("/registration", validation.Validate[dto.CreateRegistrationRequest](),
		middlewares.UserUniqueCheck[dto.CreateRegistrationRequest], handlers.CreateRegistration)
	v1.Put("/registration/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		validation.Validate[dto.EditRegistrationRequest](), handlers.EditRegistration)

	// Buy kursus
	v1.Get("/buy-kursus", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetAllBuyKursus)
	v1.Get("/buy-kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}),
		handlers.GetBuyKursus)
	v1.Post("/buy-kursus", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"siswa"}),
		validation.Validate[dto.BuyKursusRequest](),
		handlers.CreateBuyKursus)
	v1.Put("/buy-kursus/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"helpdesk"}), validation.Validate[dto.EditBuyKursus](),
		handlers.EditBuyKursus)

	// confirm-payment
	v1.Post("/confirm-payment/:id", handlers.ConfirmPayment)

	// my-course
	v1.Get("/my-course", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"siswa"}),
		handlers.GetMyCourse)
	v1.Get("/my-course/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"siswa"}),
		handlers.GetMyCourseByID)

	// Buat pertemuan
	v1.Post("/pertemuan", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"guru"}),
		validation.Validate[dto.CreatePertemuanRequest](), policy.GroupBatchOwnerPolicy("create"),
		handlers.CreatePertemuan)
	v1.Put("/pertemuan/:id", middlewares.AuthMiddleware(), middlewares.RoleAuthorization([]string{"guru"}),

		validation.Validate[dto.EditPertemuanRequest](), policy.GroupBatchOwnerPolicy("update"),
		handlers.EditPertemuan)

	// master
	v1.Get("/categories", handlers.GetAllCategories)
	v1.Get("/categories/:id",
		handlers.GetDetailCategories)

	v1.Get("/hari", handlers.GetAllHari)
	v1.Get("/hari/:id",
		handlers.GetDetailHari)

	v1.Get("/jenis-kursus", handlers.GetAllJenisKursus)
	v1.Get("/jenis-kursus/:id",
		handlers.GetDetailJenisKursus)

	v1.Get("/golongan", handlers.GetAllGolongan)
	v1.Get("/golongan/:id",
		handlers.GetDetailGolongan)

	v1.Get("/kelas", handlers.GetAllKelas)
	v1.Get("/kelas/:id",
		handlers.GetDetailKelas)

	v1.Get("/prices", handlers.GetAllPrices)
	v1.Get("/prices/:id",
		handlers.GetDetailPrice)

	v1.Get("/roles", handlers.GetAllRoles)
	v1.Get("/roles/:id",
		handlers.GetDetailRole)

	v1.Get("/status-payments", handlers.GetAllStatusPayments)
	v1.Get("/status-payments/:id",
		handlers.GetDetailStatusPayment)
}
