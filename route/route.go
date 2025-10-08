package route

import (
	"tugas-praktikum-crud/middleware"
	"tugas-praktikum-crud/app/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", service.Login)

	protected := api.Group("", middleware.AuthRequired())

	protected.Get("/profile", service.GetProfile)

	alumni := protected.Group("/alumni")
	alumni.Get("/lulusan", service.GetAlumniByTahunLulus)
	alumni.Get("/", service.GetAllAlumni)
	alumni.Get("/:id", service.GetAlumniByID)
	alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
	alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
	alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", service.GetAllPekerjaan)
	pekerjaan.Get("/trash", service.GetTrashPekerjaan)
	pekerjaan.Delete("/trash/:id", service.HardDeletePekerjaan)
    pekerjaan.Post("/trash/:id/restore", service.RestorePekerjaan)
	pekerjaan.Get("/:id", service.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
	pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
	pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
	// pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
	pekerjaan.Put("softdel/:id", service.SoftDeletePekerjaan)
}