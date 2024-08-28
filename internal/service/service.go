package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidToken = errors.New("failed to decode the provided Token")
)

func EncryptPass(pass string) string {
	h := md5.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}
