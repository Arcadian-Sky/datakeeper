package model

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrUserAuth           = errors.New("authentication failed")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("JWT is invalid")
	ErrPdataAlreatyEsists = errors.New("data already exists")
	ErrNoToken            = errors.New("no JWT")
	ErrPdataNotFound      = errors.New("data not found")

	ErrCreateBucketFailed = errors.New("failed to create bucket")
	ErrCreateBucketExists = errors.New("bucket already exists")
	ErrCreateBucketNoUser = errors.New("bucket user id not exists")
	ErrNoUserBucket       = errors.New("bucket for user id not exists")

	ErrEmptyRequestBody    = errors.New("request body is empty")
	ErrErrorRequestBody    = errors.New("failed to read request body")
	ErrFailedToDecodeJSON  = errors.New("failed to decode JSON")
	ErrStructureJSON       = errors.New("error in JSON strucrure")
	ErrInternalServer      = errors.New("internal server error")
	ErrNotAuthorized       = errors.New("user not authenticated")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidOrderNumber  = errors.New("invalid order number format")
	ErrAddOrder            = errors.New("failed to add order")
	ErrAddExistsOrder      = errors.New("order number already exists")
	ErrInvalidLoginAndPass = errors.New("invalid login/password")
	ErrLoginAlreadyTaken   = errors.New("login already taken")
	ErrEmptyResponse       = errors.New("empty response")
	ErrIncFunds            = errors.New("insufficient funds")
)

// Jtoken - JWT token
type Jtoken struct {
	Token  string
	Claims Claims
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
	Iat    int64
	Exp    int64
}

type User struct {
	ID         int64     `json:"id"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	Bucket     string    `json:"bucket"`
	LastUpdate time.Time `json:"last_update"`
}

type Data struct {
	ID       int64
	UserID   int64
	Title    string
	Type     string
	Card     string
	Login    string
	Password string
	KeyHash  string
}

type FileItem struct {
	Hash string
	Name string
	Desc string
}
