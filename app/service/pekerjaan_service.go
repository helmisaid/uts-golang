package service

import (
	"database/sql"
	"strconv"
	"strings"
	"tugas-praktikum-crud/app/model"
	"tugas-praktikum-crud/app/repository"
	"tugas-praktikum-crud/database"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaan(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")

	offset := (page - 1) * limit

	pekerjaan, err := repository.GetAllPekerjaanWithParams(database.DB, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}

	total, err := repository.CountPekerjaan(database.DB, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghitung total data pekerjaan"})
	}

	response := model.PekerjaanResponse{
		Data: pekerjaan,
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

func GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID tidak valid"})
	}
	pekerjaan, err := repository.GetPekerjaanByID(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal mengambil data pekerjaan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": pekerjaan, "message": "Data pekerjaan berhasil diambil"})
}

// GetPekerjaanByAlumniID menangani request untuk mendapatkan pekerjaan berdasarkan alumni_id.
func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Alumni ID tidak valid"})
	}
	pekerjaan, err := repository.GetPekerjaanByAlumniID(database.DB, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal mengambil data pekerjaan"})
	}
	if len(pekerjaan) == 0 {
		return c.JSON(fiber.Map{"success": true, "data": []model.PekerjaanAlumni{}, "message": "Tidak ada data pekerjaan untuk alumni ini"})
	}
	return c.JSON(fiber.Map{"success": true, "data": pekerjaan, "message": "Data pekerjaan berhasil diambil"})
}

// CreatePekerjaan menangani request untuk membuat data pekerjaan baru.
func CreatePekerjaan(c *fiber.Ctx) error {
	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Request body tidak valid"})
	}
	newPekerjaan, err := repository.CreatePekerjaan(database.DB, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal menambah pekerjaan"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": newPekerjaan, "message": "Pekerjaan berhasil ditambahkan"})
}

// UpdatePekerjaan menangani request untuk memperbarui data pekerjaan.
func UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID tidak valid"})
	}
	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "Request body tidak valid"})
	}
	updatedPekerjaan, err := repository.UpdatePekerjaan(database.DB, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan untuk diupdate"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal mengupdate pekerjaan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": updatedPekerjaan, "message": "Pekerjaan berhasil diupdate"})
}

// func DeletePekerjaan(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID tidak valid"})
// 	}
// 	err = repository.DeletePekerjaan(database.DB, id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan untuk dihapus"})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal menghapus pekerjaan"})
// 	}
// 	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus"})
// }


func SoftDeletePekerjaan(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID pekerjaan tidak valid"})
    }

    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)

    if role == "user" {
        ownerUserID, err := repository.GetPekerjaanOwnerUserID(database.DB, id)
        if err != nil {
            if err == sql.ErrNoRows {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan"})
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal memverifikasi pekerjaan"})
        }

        if ownerUserID != userID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak sesuai"})
        }
    } else if role != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Akses ditolak"})
    }
 
    err = repository.SoftDeletePekerjaan(database.DB, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal menghapus pekerjaan."})
    }

    return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus."})
}


func GetTrashPekerjaan(c *fiber.Ctx) error {
    role, ok := c.Locals("role").(string)
    if !ok {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "error":   "Akses ditolak",
        })
    }

    userID := c.Locals("user_id").(int)

    var pekerjaanList []model.PekerjaanAlumni
    var err error

    if role == "admin" {
        pekerjaanList, err = repository.GetAllTrashPekerjaanForAdmin(database.DB)
    } else if role == "user" {
        pekerjaanList, err = repository.GetAllTrashPekerjaanByUserID(database.DB, userID)
    } else {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "error":   "Akses ditolak",
        })
    }

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "error":   "Gagal mengambil data sampah dari database",
        })
    }

    if len(pekerjaanList) == 0 {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "success": true,
            "data":    []model.PekerjaanAlumni{},
            "message": "Tidak ada data pekerjaan di dalam sampah",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "success": true,
        "data":    pekerjaanList,
        "message": "Data sampah pekerjaan berhasil diambil",
    })
}

func HardDeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID pekerjaan tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "user" {
		ownerUserID, err := repository.GetPekerjaanOwnerUserID(database.DB, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan"})
		}
		if ownerUserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Akses ditolak: Anda bukan pemilik data ini"})
		}
	} else if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Akses ditolak: Peran tidak diizinkan"})
	}

	err = repository.HardDeletePekerjaan(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan di dalam trash"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal menghapus pekerjaan secara permanen"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus secara permanen"})
}

func RestorePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "ID pekerjaan tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	if role == "user" {
		ownerUserID, err := repository.GetPekerjaanOwnerUserID(database.DB, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan"})
		}
		if ownerUserID != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Akses ditolak: Anda bukan pemilik data ini"})
		}
	} else if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "Akses ditolak: Peran tidak diizinkan"})
	}

	err = repository.RestorePekerjaan(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Pekerjaan tidak ditemukan di dalam trash"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "Gagal mengembalikan pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dikembalikan"})
}