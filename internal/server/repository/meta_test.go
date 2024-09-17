package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestDataRepo_Save(t *testing.T) {

	logg := logrus.New()

	type args struct {
		ctx  context.Context
		data *model.Data
	}
	tests := []struct {
		name    string
		mock    func(mock sqlmock.Sqlmock)
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Success",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO "metadata" \(dtype, user_id, title, card_number, login, password\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
					WithArgs("type1", 1, "title1", "card1", "login1", "password1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

			},
			args: args{
				ctx: context.Background(),
				data: &model.Data{
					Type:     "type1",
					UserID:   1,
					Title:    "title1",
					Card:     "card1",
					Login:    "login1",
					Password: "password1",
				},
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "Error",
			mock: func(mock sqlmock.Sqlmock) {

				// Настроить ожидание запроса и его параметры
				mock.ExpectQuery(`INSERT INTO "metadata" \(dtype, user_id, title, card_number, login, password\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`).
					WithArgs("type1", 1, "title1", "card1", "login1", "password1").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				// WillReturnError(sql.ErrConnDone)

			},
			args: args{
				ctx: context.Background(),
				data: &model.Data{
					Type:     "type1",
					UserID:   1,
					Title:    "title1",
					Card:     "card1",
					Login:    "login1",
					Password: "password1",
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mock(mock)

			r := &DataRepo{db: db, log: logg}

			got, err := r.Save(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DataRepo.Save() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataRepo_GetList(t *testing.T) {
	logg := logrus.New()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		mock    func(mock sqlmock.Sqlmock)
		args    args
		want    []model.Data
		wantErr bool
	}{
		{
			name: "Success",
			mock: func(mock sqlmock.Sqlmock) {
				// Настройка мока для успешного выполнения запроса
				rows := sqlmock.NewRows([]string{"id", "dtype", "title", "card_number", "login", "password"}).
					AddRow(1, "type1", "title1", "card1", "login1", "password1").
					AddRow(2, "type2", "title2", "card2", "login2", "password2")

				mock.ExpectQuery(`SELECT id, dtype, title, card_number, login, password FROM metadata WHERE user_id = \$1 ORDER BY id`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: 1,
				},
			},
			want: []model.Data{
				{
					ID:       1,
					Type:     "type1",
					Title:    "title1",
					Card:     "card1",
					Login:    "login1",
					Password: "password1",
				},
				{
					ID:       2,
					Type:     "type2",
					Title:    "title2",
					Card:     "card2",
					Login:    "login2",
					Password: "password2",
				},
			},
			wantErr: false,
		},
		{
			name: "QueryError",
			mock: func(mock sqlmock.Sqlmock) {
				// Настройка мока для вызова ошибки запроса
				mock.ExpectQuery(`SELECT id, dtype, title, card_number, login, password FROM metadata WHERE user_id = \$1 ORDER BY id`).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ScanError",
			mock: func(mock sqlmock.Sqlmock) {
				// Настройка мока для успешного выполнения запроса, но с ошибкой сканирования
				rows := sqlmock.NewRows([]string{"id", "dtype", "title", "card_number", "login", "password"}).
					AddRow("wrong_type", "type1", "title1", "card1", "login1", "password1") // Wrong type for `id`

				mock.ExpectQuery(`SELECT id, dtype, title, card_number, login, password FROM metadata WHERE user_id = \$1 ORDER BY id`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "IterationError",
			mock: func(mock sqlmock.Sqlmock) {
				// Настройка мока для успешного выполнения запроса, но с ошибкой итерации
				rows := sqlmock.NewRows([]string{"id", "dtype", "title", "card_number", "login", "password"}).
					AddRow(1, "type1", "title1", "card1", "login1", "password1")

				mock.ExpectQuery(`SELECT id, dtype, title, card_number, login, password FROM metadata WHERE user_id = \$1 ORDER BY id`).
					WithArgs(1).
					WillReturnRows(rows)

				rows.RowError(0, fmt.Errorf("row iteration error"))
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: 1,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mock(mock)

			r := &DataRepo{db: db, log: logg}

			got, err := r.GetList(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataRepo.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.ElementsMatch(t, tt.want, got, "DataRepo.GetList() = %v, want %v", got, tt.want)
		})
	}
}
