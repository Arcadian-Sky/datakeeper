package jwtrule

import (
	"context"
	"fmt"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/golang-jwt/jwt"
)

type ctxKey string

var (
	CtxKeyUserID ctxKey = "userID"
)

// Generate generates new JWT token
func Generate(userid int64, key string) (model.Jtoken, error) {
	now := time.Now()
	claims := model.Claims{UserID: userid, Iat: now.Unix(),
		Exp: now.Add(time.Minute * 60).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  claims.UserID,
		"iat": claims.Iat,
		"exp": claims.Exp,
	})
	tokenString, err := token.SignedString([]byte(key))
	return model.Jtoken{Claims: claims, Token: tokenString}, err
}

// Validate checks JWT token and converts to structs.Jtoken
func Validate(tokenString string, key string) (model.Jtoken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})

	if claimsMap, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims := model.Claims{
			UserID: int64(claimsMap["id"].(float64)),
			Iat:    int64(claimsMap["iat"].(float64)),
			Exp:    int64(claimsMap["exp"].(float64)),
		}
		jtoken := model.Jtoken{Token: tokenString, Claims: claims}

		return jtoken, nil
	}

	return model.Jtoken{}, err

}

// SetUserIDToCTX add userID to the context.
func SetUserIDToCTX(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, value)
}

// Получаем значение из контекста
func GetUserIDFromCTX(ctx context.Context) int64 {
	if iUserID, ok := ctx.Value(CtxKeyUserID).(int); ok {
		return int64(iUserID)
	}
	return -1
}
