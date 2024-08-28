package model

import (
	"errors"

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
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Bucket   string `json:"bucket"`
}

// type Data struct {
// 	UserID int64  `json:"id"`
// 	Bucket string `json:"bucket"`
// 	Data   string `json:"data"`
// }

type Data struct {
	ID          int64
	Name        string
	Type        string
	KeyHash     string
	PrivateData string
}
