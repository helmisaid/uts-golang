package service

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"tugas-praktikum-crud/app/model"
	"tugas-praktikum-crud/app/repository"
	"tugas-praktikum-crud/database"

	"github.com/gofiber/fiber/v2"
)

func GetAllAlumni(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")

	offset := (page - 1) * limit

	alumni, err := repository.GetAllAlumniWithParams(database.DB, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data alumni"})
	}

	total, err := repository.CountAlumni(database.DB, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghitung total data alumni"})
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit, 
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

func GetAllAlumniSimple(c *fiber.Ctx) error {
    alumni, err := repository.GetAllAlumni(database.DB)
    if err != nil {
        log.Println("Error getting all alumni:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Gagal mengambil data alumni",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    alumni,
        "message": "Semua data alumni berhasil diambil",
    })
}

func GetAllTrashAlumni(c *fiber.Ctx) error {
    alumni, err := repository.GetAllTrashAlumni(database.DB)
    if err != nil {
        log.Println("Error getting all alumni:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Gagal mengambil data alumni",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    alumni,
        "message": "Semua data alumni berhasil diambil",
    })
}

func GetAlumniByTahunLulus(c *fiber.Ctx) error {
	tahunLulusStr := c.Query("tahun_lulus")

	if tahunLulusStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Parameter tahun_lulus wajib diisi",
		})
	}

	tahunLulus, err := strconv.Atoi(tahunLulusStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Tahun lulus tidak valid",
		})
	}

	alumni, err := repository.GetAlumniByTahunLulus(database.DB, tahunLulus)
	if err != nil {
		log.Println("Error getting alumni by Tahun Lulus:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data alumni",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    alumni,
		"message": "Data alumni berhasil diambil",
	})
}

// GetAlumniByID menangani request untuk mendapatkan alumni berdasarkan ID.
func GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "ID tidak valid",
		})
	}

	alumni, err := repository.GetAlumniByID(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Alumni tidak ditemukan",
			})
		}
		log.Println("Error getting alumni by ID:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data alumni",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    alumni,
		"message": "Data alumni berhasil diambil",
	})
}

// CreateAlumni menangani request untuk membuat alumni baru.
func CreateAlumni(c *fiber.Ctx) error {
	var req model.Alumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Request body tidak valid",
		})
	}

	if req.NIM == "" || req.Nama == "" || req.Email == "" || req.Jurusan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Field nim, nama, email, dan jurusan harus diisi",
		})
	}

	newAlumni, err := repository.CreateAlumni(database.DB, req)
	if err != nil {
		log.Println("Error creating alumni:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menambah alumni. Pastikan NIM dan email belum digunakan",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newAlumni,
		"message": "Alumni berhasil ditambahkan",
	})
}

// UpdateAlumni menangani request untuk memperbarui data alumni.
func UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID tidak valid"})
	}

	var req model.Alumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Request body tidak valid"})
	}

	updatedAlumni, err := repository.UpdateAlumni(database.DB, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Alumni tidak ditemukan untuk diupdate",
			})
		}
		log.Println("Error updating alumni:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengupdate alumni",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedAlumni,
		"message": "Alumni berhasil diupdate",
	})
}

// DeleteAlumni menangani request untuk menghapus data alumni.
func DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID tidak valid"})
	}

	err = repository.DeleteAlumni(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Alumni tidak ditemukan untuk dihapus",
			})
		}
		log.Println("Error deleting alumni:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menghapus alumni",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}


