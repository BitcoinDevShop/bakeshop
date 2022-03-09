package bakedgood

import (
	"database/sql"
	"errors"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS bakedgoods(
		_id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        name TEXT NOT NULL UNIQUE,
        macaroon TEXT NOT NULL,
        status TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS payments(
		_id INTEGER PRIMARY KEY AUTOINCREMENT,
        macaroon_id TEXT NOT NULL,
        amount INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
    `

	_, err := r.db.Exec(query)
	return err
}
