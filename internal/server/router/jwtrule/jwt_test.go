package jwtrule

import (
	"context"
	"testing"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	type args struct {
		userid int64
		key    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful token generation",
			args: args{
				userid: 12345,
				key:    "test-secret-key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.userid, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify the token
			token, err := jwt.Parse(got.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.args.key), nil
			})
			if err != nil {
				t.Errorf("jwt.Parse() error = %v", err)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				t.Errorf("jwt token is not valid")
				return
			}

			// Verify claims
			assert.Equal(t, tt.args.userid, int64(claims["id"].(float64)), "Expected userID to be equal")
			assert.WithinDuration(t, time.Unix(int64(claims["iat"].(float64)), 0), time.Now(), time.Minute, "Expected iat claim to be within 1 minute of current time")
			assert.WithinDuration(t, time.Unix(int64(claims["exp"].(float64)), 0), time.Now().Add(time.Minute*60), time.Minute, "Expected exp claim to be within 1 minute of 60 minutes from now")
		})
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		tokenString string
		key         string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Jtoken
		wantErr bool
	}{
		{
			name: "successful validation",
			args: args{
				tokenString: generateTestToken(12345, "test-secret-key"),
				key:         "test-secret-key",
			},
			want: model.Jtoken{
				Claims: model.Claims{
					UserID: 12345,
					Iat:    time.Now().Unix(),
					Exp:    time.Now().Add(time.Minute * 60).Unix(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid token",
			args: args{
				tokenString: "invalid-token",
				key:         "test-secret-key",
			},
			want:    model.Jtoken{},
			wantErr: true,
		},
		{
			name: "invalid signing method",
			args: args{
				tokenString: generateTestToken(12345, "wrong-key"),
				key:         "test-secret-key",
			},
			want:    model.Jtoken{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Validate(tt.args.tokenString, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !assert.Equal(t, tt.want, got) {
			// 	t.Errorf("Validate() = %v, want %v", got, tt.want)
			// }
		})
	}
}

// Helper function to generate a test JWT token
func generateTestToken(userID int64, key string) string {
	now := time.Now()
	claims := model.Claims{
		UserID: userID,
		Iat:    now.Unix(),
		Exp:    now.Add(time.Minute * 60).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  claims.UserID,
		"iat": claims.Iat,
		"exp": claims.Exp,
	})
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}

// TestSetUserIDToCTX tests the SetUserIDToCTX function
func TestSetUserIDToCTX(t *testing.T) {
	ctx := context.Background()
	userID := int64(12345)

	// Set user ID to context
	newCtx := SetUserIDToCTX(ctx, int(userID))

	// Get the user ID from the new context
	got := GetUserIDFromCTX(newCtx)

	assert.Equal(t, userID, got, "UserID should be set and retrieved correctly")
}

// TestGetUserIDFromCTX tests the GetUserIDFromCTX function
func TestGetUserIDFromCTX(t *testing.T) {
	ctx := context.Background()

	// Test with a context that has no userID
	got := GetUserIDFromCTX(ctx)
	assert.Equal(t, int64(-1), got, "UserID should return -1 when not set")

	// Set a userID in the context
	userID := int64(12345)
	ctxWithUserID := SetUserIDToCTX(ctx, int(userID))

	// Test with a context that has a userID
	got = GetUserIDFromCTX(ctxWithUserID)
	assert.Equal(t, userID, got, "UserID should be retrieved correctly from context")
}
