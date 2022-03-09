package bakedgood

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

// TODO this could probably just be MacaroonDetails
type BakedGood struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Macaroon  string    `json:"macaroon"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *SQLiteRepository) CreateBakedGood(bg BakedGood) (*BakedGood, error) {
	_, err := r.db.Exec("INSERT INTO bakedgoods(uuid, name, macaroon, status) values(?,?,?,?)", bg.Id, bg.Name, bg.Macaroon, bg.Status)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	// TODO we could use sqlite's integer id but we don't
	// _id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &bg, nil
}

func (r *SQLiteRepository) AllBakedGoods() ([]BakedGood, error) {
	rows, err := r.db.Query("SELECT * FROM bakedgoods")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []BakedGood
	for rows.Next() {
		// We don't use this but we need to read it
		var id int64
		var bg BakedGood
		if err := rows.Scan(&id, &bg.Id, &bg.Name, &bg.Macaroon, &bg.Status, &bg.CreatedAt); err != nil {
			return nil, err
		}
		all = append(all, bg)
	}
	return all, nil
}

func (r *SQLiteRepository) GetBakedGoodByUuid(uuid string) (*BakedGood, error) {
	row := r.db.QueryRow("SELECT * FROM bakedgoods WHERE uuid = ?", uuid)

	var id int64
	var bg BakedGood
	if err := row.Scan(&id, &bg.Id, &bg.Name, &bg.Macaroon, &bg.Status, &bg.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &bg, nil
}

func (r *SQLiteRepository) UpdateBakedGood(uuid string, updated BakedGood) (*BakedGood, error) {
	if uuid == "" {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE bakedgoods SET uuid = ?, name = ?, macaroon = ?, status = ? WHERE uuid = ?", updated.Id, updated.Name, updated.Macaroon, updated.Status, uuid)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) DeleteBakedGood(uuid string) error {
	res, err := r.db.Exec("DELETE FROM bakedgoods WHERE uuid = ?", uuid)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
