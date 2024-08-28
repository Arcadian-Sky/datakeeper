package repository

import (
	"context"
	"database/sql"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/sirupsen/logrus"
)

type DataRepository interface {
	Add(ctx context.Context, data model.Data) (int64, error)
	GetList(ctx context.Context, data model.Data) (int64, error)
}

type dataRepo struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewDataRepository(dbd *sql.DB, lg *logrus.Logger) *dataRepo {
	p := &dataRepo{
		db:  dbd,
		log: lg,
	}
	return p
}

func (d *dataRepo) Add(ctx context.Context, data model.Data) (int64, error) {
	return 0, nil
}

func (d *dataRepo) GetList(ctx context.Context, data model.Data) (int64, error) {
	return 0, nil
}
