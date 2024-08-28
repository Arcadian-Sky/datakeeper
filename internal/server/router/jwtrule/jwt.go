package jwtrule

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/golang-jwt/jwt"
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

type ctxkey string

var (
	userID ctxkey = "userID"
)

// SetUserIDToCTX add userID to the context.
func SetUserIDToCTX(ctx context.Context, value int) context.Context {
	return context.WithValue(ctx, userID, value)
}

func SetUserIDFromCTX(ctx context.Context) int64 {
	// Получаем значение из контекста
	if strUserID, ok := ctx.Value(userID).(string); ok {
		// Преобразуем строку в int64
		userID, err := strconv.ParseInt(strUserID, 10, 64)
		if err != nil {
			return userID
		}
	}
	return 0
}
