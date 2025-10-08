package repository

import (
	"database/sql"
	"fmt"
	"time"
	"tugas-praktikum-crud/app/model"
)

func CreateAlumni(db *sql.DB, req model.Alumni) (*model.Alumni, error) {
	var id int
	query := `
        INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
        RETURNING id`

	err := db.QueryRow(
		query,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now(),
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return GetAlumniByID(db, id)
}

func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
	var alumniList []model.Alumni

	query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at, user_id 
        FROM alumni 
        ORDER BY created_at DESC` 

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
			&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
			&alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}

	return alumniList, nil
}

func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	alumni := new(model.Alumni)

	query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at 
        FROM alumni 
        WHERE id = $1`

	err := db.QueryRow(query, id).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
		&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func GetAlumniByTahunLulus(db *sql.DB, tahun_lulus int) ([]model.Alumni, error) {
	var alumniList []model.Alumni

	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at 
			FROM alumni 
			WHERE tahun_lulus = $1`
	
	rows, err := db.Query(query, tahun_lulus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
			&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
			&alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}

	return alumniList, nil
}


// func GetAlumniByStatusPekerjaan(db *sql.DB, id int) (*model.Alumni, error) {
// 	alumni := new(model.Alumni)

// 	query := `
//         SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
//                no_telepon, alamat, created_at, updated_at 
//         FROM alumni 
//         WHERE id = $1`

// 	err := db.QueryRow(query, id).Scan(
// 		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
// 		&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
// 		&alumni.CreatedAt, &alumni.UpdatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return alumni, nil
// }



func UpdateAlumni(db *sql.DB, id int, req model.Alumni) (*model.Alumni, error) {
	query := `
        UPDATE alumni 
        SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, 
            email = $6, no_telepon = $7, alamat = $8, updated_at = $9 
        WHERE id = $10`

	result, err := db.Exec(
		query,
		req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, time.Now(), id,
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

	return GetAlumniByID(db, id)
}

func DeleteAlumni(db *sql.DB, id int) error {
	query := `DELETE FROM alumni WHERE id = $1`

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

func GetAllAlumniWithParams(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	var alumniList []model.Alumni

	allowedSortBy := map[string]string{
		"id":         "id",
		"nama":       "nama",
		"nim":        "nim",
		"angkatan":   "angkatan",
		"created_at": "created_at",
	}
	sortColumn, ok := allowedSortBy[sortBy]
	if !ok {
		sortColumn = "id" 
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at 
        FROM alumni 
        WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3`, sortColumn, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
			&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
			&alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}
	return alumniList, nil
}

func CountAlumni(db *sql.DB, search string) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1`
	err := db.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAllTrashAlumni(db *sql.DB) ([]model.Alumni, error) {
	var alumniList []model.Alumni

	query := `
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
               no_telepon, alamat, created_at, updated_at, user_id 
        FROM alumni
		WHERE deleted_at IS NOT NULL 
        `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
			&alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon, &alumni.Alamat,
			&alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}

	return alumniList, nil
}