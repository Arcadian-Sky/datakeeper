package repository

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepo_Register(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *model.User
	}
	logg := logrus.New()
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
		mock    func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Successful Registration",
			args: args{
				ctx:  context.Background(),
				user: &model.User{Login: "newuser", Password: "password123"},
			},
			want:    1,
			wantErr: false,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE login = \$1`).
					WithArgs("newuser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

					// WillReturnError(sql.ErrNoRows)

				mock.ExpectQuery(`INSERT INTO "user" \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("newuser", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))
			},
		},
		{
			name: "User Already Exists",
			args: args{
				ctx:  context.Background(),
				user: &model.User{Login: "existinguser", Password: "password123"},
			},
			want:    0,
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE login = \$1`).
					WithArgs("existinguser").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
		{
			name: "Failed Password Hash",
			args: args{
				ctx:  context.Background(),
				user: &model.User{Login: "newuser", Password: ""},
			},
			want:    0,
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id FROM "user" WHERE login = \$1`).
					WithArgs("newuser").
					WillReturnError(sql.ErrNoRows)
			},
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

			r := &UserRepo{db: db, log: logg}
			got, err := r.Register(tt.args.ctx, tt.args.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got != tt.want {
					t.Errorf("UserRepo.Register() = %v, want %v", got, tt.want)
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestUserRepo_Auth(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *model.User
	}
	logg := logrus.New()
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
		mock    func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Successful Authentication",
			args: args{
				ctx:  context.Background(),
				user: &model.User{Login: "existinguser", Password: "password123"},
			},
			want:    &model.User{ID: 1, Login: "existinguser"},
			wantErr: false,
			mock: func(mock sqlmock.Sqlmock) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mock.ExpectQuery(`SELECT id, password FROM "user" WHERE login=\$1`).
					WithArgs("existinguser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword))
			},
		},
		{
			name: "Invalid Credentials",
			args: args{
				ctx:  context.Background(),
				user: &model.User{Login: "wronguser", Password: "wrongpass"},
			},
			want:    &model.User{},
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, password FROM "user" WHERE login=\$1`).
					WithArgs("wronguser").
					WillReturnError(sql.ErrNoRows)
			},
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

			r := &UserRepo{db: db, log: logg}
			got, err := r.Auth(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.ID, tt.want.ID) {
					t.Errorf("UserRepo.Auth() = %v, want %v", got.ID, tt.want.ID)
				}
				if !reflect.DeepEqual(got.Login, tt.want.Login) {
					t.Errorf("UserRepo.Auth() = %v, want %v", got.Login, tt.want.Login)
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations: %s", err)
				}
			}
		})
	}
}

func TestUserRepo_SetLastUpdate(t *testing.T) {
	type args struct {
		ctx  context.Context
		user *model.User
	}
	tm := time.Now()
	logg := logrus.New()
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
		mock    func(mock sqlmock.Sqlmock)
	}{
		{
			name: "Successful Last Update",
			args: args{
				ctx:  context.Background(),
				user: &model.User{ID: 1, LastUpdate: tm},
			},
			want:    &model.User{ID: 1, LastUpdate: tm},
			wantErr: false,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "user" SET last_update = \$1 WHERE id = \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Failed Last Update",
			args: args{
				ctx:  context.Background(),
				user: &model.User{ID: 1, LastUpdate: tm},
			},
			want:    &model.User{},
			wantErr: true,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE "user" SET last_update = \$1 WHERE id = \$2`).
					WithArgs(sqlmock.AnyArg(), 1).
					WillReturnError(sql.ErrNoRows)
			},
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

			r := &UserRepo{db: db, log: logg}
			got, err := r.SetLastUpdate(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.SetLastUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UserRepo.SetLastUpdate() = %v, want %v", got, tt.want)
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("there were unfulfilled expectations")
				}
			}
		})
	}
}
