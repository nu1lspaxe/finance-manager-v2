package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey []byte
}

func NewJWTManager(secretKey []byte) *JWTManager {
	return &JWTManager{secretKey}
}

func (j *JWTManager) Clone() *JWTManager {
	return &JWTManager{j.secretKey}
}

type FMJWTClaims struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTManager) Verify(tokenString string) (*FMJWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&FMJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, NewRecordError(ErrTokenAlg)
			}
			return j.secretKey, nil

		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*FMJWTClaims)
	if !ok {
		return nil, NewRecordError(ErrTokenInvalid)
	}

	audienceValid := false
	for _, aud := range claims.Audience {
		if aud == "records" {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return nil, NewRecordError(ErrTokenInvalid)
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, NewRecordError(ErrTokenInvalid)
	}

	return claims, nil
}
