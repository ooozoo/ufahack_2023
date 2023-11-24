package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"ufahack_2023/internal/domain"
)

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrFailedToParseClaims     = errors.New("failed to parse claims")
)

type Claims struct {
	UID domain.ID
	Exp time.Time
}

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

func ParseToken(token string, secret string) (*Claims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrFailedToParseClaims
	}

	uidRaw, ok := claims["uid"]
	if !ok {
		return nil, ErrFailedToParseClaims
	}

	uidStr, ok := uidRaw.(string)
	if !ok {
		return nil, ErrFailedToParseClaims
	}

	uid, err := uuid.Parse(uidStr)
	if !ok {
		return nil, fmt.Errorf("%w: %w", ErrFailedToParseClaims, err)
	}

	expRaw, ok := claims["exp"]
	if !ok {
		return nil, ErrFailedToParseClaims
	}

	expFloat, ok := expRaw.(float64)
	if !ok {
		return nil, ErrFailedToParseClaims
	}

	exp := time.Unix(int64(expFloat), 0)

	return &Claims{
		UID: uid,
		Exp: exp,
	}, nil
}
