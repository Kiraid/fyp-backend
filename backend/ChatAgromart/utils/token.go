package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func VerifyToken(token string) (int64, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, errors.New("couldnt parse token")
	}
	tokenisValid := parsed.Valid
	if !tokenisValid {
		return 0, errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	userId := int64(claims["userId"].(float64))
	return userId, nil
}
