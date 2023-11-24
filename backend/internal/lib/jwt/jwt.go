package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"ufahack_2023/internal/domain"
)

func NewToken(user *domain.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
