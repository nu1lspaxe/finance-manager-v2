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

type UserJWTClaims struct {
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTManager) Verify(tokenString string) (*UserJWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, NewBankRecordError(ErrTokenAlg)
			}
			return j.secretKey, nil

		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserJWTClaims)
	if !ok {
		return nil, NewBankRecordError(ErrTokenInvalid)
	}

	audienceValid := false
	for _, aud := range claims.Audience {
		if aud == "records_bank" {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return nil, NewBankRecordError(ErrTokenInvalid)
	}

	if claims.ExpiresAt.Time.Unix() < time.Now().Unix() {
		return nil, NewBankRecordError(ErrTokenInvalid)
	}

	return claims, nil
}
