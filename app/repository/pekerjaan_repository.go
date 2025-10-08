package repository

import (
	"database/sql"
	"fmt"
	"time"
	"tugas-praktikum-crud/app/model"
)

func CreatePekerjaan(db *sql.DB, req model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	var id int
	query := `
        INSERT INTO pekerjaan_alumni (
            alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
            gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
            deskripsi_pekerjaan, created_at, updated_at
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
        RETURNING id`

	err := db.QueryRow(
		query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, time.Now(), time.Now(),
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return GetPekerjaanByID(db, id)
}

func GetAllPekerjaan(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	var pekerjaanList []model.PekerjaanAlumni
	query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
               gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
               deskripsi_pekerjaan, created_at, updated_at 
        FROM pekerjaan_alumni 
        ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
			&p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func GetPekerjaanByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
               gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
               deskripsi_pekerjaan, created_at, updated_at 
        FROM pekerjaan_alumni 
        WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
		&p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPekerjaanByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	var pekerjaanList []model.PekerjaanAlumni
	query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
               gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
               deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni 
        WHERE alumni_id = $1 AND is_deleted IS NULL
        ORDER BY tanggal_mulai_kerja DESC`

	rows, err := db.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
			&p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func UpdatePekerjaan(db *sql.DB, id int, req model.PekerjaanAlumni) (*model.PekerjaanAlumni, error) {
	query := `
        UPDATE pekerjaan_alumni 
        SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, 
            lokasi_kerja = $5, gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, 
            status_pekerjaan = $9, deskripsi_pekerjaan = $10, updated_at = $11 
        WHERE id = $12`

	result, err := db.Exec(
		query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, req.TanggalMulaiKerja, req.TanggalSelesaiKerja, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, time.Now(), id,
	)

	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return GetPekerjaanByID(db, id)
}

// func DeletePekerjaan(db *sql.DB, id int) error {
// 	query := "DELETE FROM pekerjaan_alumni WHERE id = $1"
// 	result, err := db.Exec(query, id)
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return sql.ErrNoRows
// 	}
// 	return nil
// }

func SoftDeletePekerjaan(db *sql.DB, id int) error {
	query := "UPDATE pekerjaan_alumni SET is_deleted = NOW() WHERE id = $1 AND is_deleted IS NULL"
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetPekerjaanOwnerUserID(db *sql.DB, pekerjaanID int) (int, error) {
    var userID int
    query := `
        SELECT a.user_id 
        FROM pekerjaan_alumni pa
        JOIN alumni a ON pa.alumni_id = a.id
        WHERE pa.id = $1`
    
    err := db.QueryRow(query, pekerjaanID).Scan(&userID)
    return userID, err
}

func GetAllPekerjaanWithParams(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	var pekerjaanList []model.PekerjaanAlumni

	allowedSortBy := map[string]string{
		"id":              "id",
		"nama_perusahaan": "nama_perusahaan",
		"posisi_jabatan":  "posisi_jabatan",
		"lokasi_kerja":    "lokasi_kerja",
		"tanggal_mulai_kerja": "tanggal_mulai_kerja",
	}
	sortColumn, ok := allowedSortBy[sortBy]
	if !ok {
		sortColumn = "id"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query := fmt.Sprintf(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
               gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
               deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni
        WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR lokasi_kerja ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3`, sortColumn, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func CountPekerjaan(db *sql.DB, search string) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR lokasi_kerja ILIKE $1`
	err := db.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAllTrashPekerjaanByUserID(db *sql.DB, userID int) ([]model.PekerjaanAlumni, error) {
    var pekerjaanList []model.PekerjaanAlumni

    query := `
        SELECT 
            pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri, 
            pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
            pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at, pa.is_deleted
        FROM 
            pekerjaan_alumni pa
        JOIN 
            alumni a ON pa.alumni_id = a.id
        WHERE 
            a.user_id = $1 AND pa.is_deleted IS NOT NULL
        ORDER BY 
            pa.is_deleted DESC`

    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("gagal menjalankan query: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var p model.PekerjaanAlumni
        if err := rows.Scan(
            &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
            &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
            &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted,
        ); err != nil {
            return nil, fmt.Errorf("gagal memindai baris data pekerjaan: %w", err)
        }
        pekerjaanList = append(pekerjaanList, p)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error pada rows setelah iterasi: %w", err)
    }
    
    return pekerjaanList, nil
}

func GetAllTrashPekerjaanForAdmin(db *sql.DB) ([]model.PekerjaanAlumni, error) {
    var pekerjaanList []model.PekerjaanAlumni

    query := `
        SELECT 
            id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
            lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
            status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_deleted
        FROM 
            pekerjaan_alumni
        WHERE 
            is_deleted IS NOT NULL
        ORDER BY 
            is_deleted DESC`

    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("gagal menjalankan query admin: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var p model.PekerjaanAlumni
        if err := rows.Scan(
            &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
            &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
            &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted,
        ); err != nil {
            return nil, fmt.Errorf("gagal memindai baris data pekerjaan admin: %w", err)
        }
        pekerjaanList = append(pekerjaanList, p)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error pada rows setelah iterasi admin: %w", err)
    }
    
    return pekerjaanList, nil
}

func HardDeletePekerjaan(db *sql.DB, id int) error {
    query := "DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_deleted IS NOT NULL"
    
    result, err := db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("gagal mengeksekusi hard delete: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("gagal mendapatkan baris yang terpengaruh: %w", err)
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }

    return nil
}

func RestorePekerjaan(db *sql.DB, id int) error {
    query := "UPDATE pekerjaan_alumni SET is_deleted = NULL WHERE id = $1 AND is_deleted IS NOT NULL"

    result, err := db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("gagal mengeksekusi restore: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("gagal mendapatkan baris yang terpengaruh setelah restore: %w", err)
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}