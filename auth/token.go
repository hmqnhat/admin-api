package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	Issuer = "aria-inc"
)

type Claims struct {
	jwt.RegisteredClaims
}

func NewClaims(accountID string, now time.Time) Claims {
	AccessExpiration, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_TIME"))
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(AccessExpiration) * time.Minute)),
			Issuer:    Issuer,
			Subject:   accountID,
			ID:        uuid.NewString(),
		},
	}
}

func Sign(c *Claims) (string, error) {
	key := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(key))
}

func UnSign(signedToken string) (*Claims, error) {
	key := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "access token parse failure")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.Wrap(jwt.ErrTokenInvalidClaims, "invalid token")
}

func IsExpired(expireTime int64) bool {
	return time.Now().Unix() > expireTime
}

func GenerateRefreshToken(now time.Time) string {
	RefreshExpiration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP_TIME"))
	return fmt.Sprintf("%s_%d", uuid.NewString(), now.Add(time.Duration(RefreshExpiration)*time.Hour).Unix())
}
