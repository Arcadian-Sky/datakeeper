package repository

import (
	"context"
	"database/sql"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/sirupsen/logrus"
)

var (
	DataType_CARD    = "CARD"
	DataType_LOGPASS = "LOGPASS"
)

type DataRepository interface {
	Save(ctx context.Context, data *model.Data) (int64, error)
	GetList(ctx context.Context, user *model.User) ([]model.Data, error)
}

type DataRepo struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewDataRepository(dbd *sql.DB, lg *logrus.Logger) *DataRepo {
	p := &DataRepo{
		db:  dbd,
		log: lg,
	}
	return p
}

func (d *DataRepo) Save(ctx context.Context, data *model.Data) (int64, error) {
	// Insert new user
	insertQuery := `INSERT INTO "metadata" (dtype, user_id, title, card_number, login, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := d.db.QueryRowContext(ctx, insertQuery, data.Type, data.UserID, data.Title, data.Card, data.Login, data.Password).Scan(&data.ID)
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (d *DataRepo) GetList(ctx context.Context, user *model.User) ([]model.Data, error) {
	query := `SELECT id, dtype, title, card_number, login, password FROM metadata WHERE user_id = $1 ORDER BY id`
	rows, err := d.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		d.log.WithError(err).Error("Failed to get metadata")
		return nil, err
	}
	defer rows.Close()

	var datalist []model.Data
	for rows.Next() {
		var data model.Data
		if err := rows.Scan(&data.ID, &data.Type, &data.Title, &data.Card, &data.Login, &data.Password); err != nil {
			d.log.WithError(err).Error("Failed to scan data")
			return nil, err
		}
		datalist = append(datalist, data)
	}

	if err := rows.Err(); err != nil {
		d.log.WithError(err).Error("Error while iterating rows")
		return nil, err
	}

	return datalist, nil
}
