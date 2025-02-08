package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"
const refreshSecretKey = "refreshsecret"

func GenerateToken(userName string, userID int64, role string) (string,string, error){
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":userName,
		"userId":userID,
		"role": role,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})
	access, err := accesstoken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	})

	refresh, err := refreshToken.SignedString([]byte(refreshSecretKey))
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func VerifyToken(token string) (int64, string, error){
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil{
		return 0, "", errors.New("couldnt parse token")
	}
	tokenisValid := parsed.Valid
	if !tokenisValid{
		return 0, "", errors.New("INVALID TOKEN")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("INVALID TOKEN CLAIMS")
	}

	userId := int64(claims["userId"].(float64))
	role := claims["role"].(string)

	return userId, role,  nil
}

func VerifyRefreshToken(token string) (int64, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(refreshSecretKey), nil
	})
	if err != nil {
		return 0, errors.New("couldn't parse token")
	}
	if !parsed.Valid {
		return 0, errors.New("INVALID TOKEN")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("INVALID TOKEN CLAIMS")
	}

	userId := int64(claims["userId"].(float64))
	return userId, nil
}

