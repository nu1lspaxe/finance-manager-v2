package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func GetJWTToken(userId int64, issueTime time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"iat":     issueTime.Unix(),
		"exp":     issueTime.Add(TOKEN_EXPIRY).Unix(),
	})
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWTToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, err
	}
	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, err
	}
	return int64(userId), nil
}
