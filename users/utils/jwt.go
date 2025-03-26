package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJWTManager(secretKey []byte, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (j *JWTManager) Clone() *JWTManager {
	return &JWTManager{j.secretKey, j.tokenDuration}
}

type UserJWTClaims struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTManager) Generate(userId int64, issueTime time.Time) (string, error) {
	claims := UserJWTClaims{
		UserId: userId,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issueTime),
			ExpiresAt: jwt.NewNumericDate(issueTime.Add(j.tokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTManager) Verify(tokenString string) (*UserJWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, NewUserError(ErrTokenAlg)
			}
			return j.secretKey, nil

		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserJWTClaims)
	if !ok {
		return nil, NewUserError(ErrTokenInvalid)
	}

	return claims, nil
}
