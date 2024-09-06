package repository

import (
	"context"
	"database/sql"

	pb "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(ctx context.Context, user *model.User) (int64, error)
	Auth(ctx context.Context, user *model.User) (*model.User, error)
	SetLastUpdate(ctx context.Context, user *model.User) (*model.User, error)
}

type UserRepo struct {
	pb.UnimplementedUserServiceServer
	db  *sql.DB
	log *logrus.Logger
}

func NewUserRepository(dbsql *sql.DB, logger *logrus.Logger) *UserRepo {
	p := &UserRepo{
		db:  dbsql,
		log: logger,
	}
	return p
}

func (r *UserRepo) Register(ctx context.Context, user *model.User) (int64, error) {
	var existingID int64
	query := `SELECT id FROM "user" WHERE login = $1`
	err := r.db.QueryRowContext(ctx, query, user.Login).Scan(&existingID)
	if err == nil {
		r.log.WithError(model.ErrLoginAlreadyTaken).Warning(model.ErrLoginAlreadyTaken.Error())
		return 0, model.ErrLoginAlreadyTaken
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Insert new user
	insertQuery := `INSERT INTO "user" (login, password) VALUES ($1, $2) RETURNING id`
	var userID int64
	err = r.db.QueryRowContext(ctx, insertQuery, user.Login, hashedPassword).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepo) Auth(ctx context.Context, user *model.User) (*model.User, error) {
	var storedUser model.User

	query := `SELECT id, password FROM "user" WHERE login=$1`
	err := r.db.QueryRowContext(ctx, query, user.Login).Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return &storedUser, model.ErrInvalidLoginAndPass
		}
		return &storedUser, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return &storedUser, model.ErrInvalidLoginAndPass
	}
	storedUser.Login = user.Login
	return &storedUser, nil
}

func (r *UserRepo) SetLastUpdate(ctx context.Context, user *model.User) (*model.User, error) {
	// Insert new user
	insertQuery := `UPDATE "user" SET last_update = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, insertQuery, user.LastUpdate, user.ID)
	if err != nil {
		return user, err
	}

	return user, nil
}
