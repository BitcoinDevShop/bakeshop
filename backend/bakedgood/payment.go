package bakedgood

import (
	"errors"
	"strconv"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Payment struct {
	Id         string    `json:"id"`
	MacaroonId string    `json:"macaroon_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

func (r *SQLiteRepository) CreatePayment(p Payment) (*Payment, error) {
	res, err := r.db.Exec("INSERT INTO payments(macaroon_id, amount) values(?,?)", p.MacaroonId, p.Amount)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	sqlId, err := res.LastInsertId()
	p.Id = strconv.FormatInt(sqlId, 10)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *SQLiteRepository) GetPaymentsByMacaroonId(macaroonId string) ([]Payment, error) {
	rows, err := r.db.Query("SELECT * FROM payments WHERE macaroon_id = ? ORDER BY created_at", macaroonId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.Id, &p.MacaroonId, &p.Amount, &p.CreatedAt); err != nil {
			return nil, err
		}
		all = append(all, p)
	}
	return all, nil
}
