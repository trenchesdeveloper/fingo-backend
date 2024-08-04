package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type JWTToken struct {
	config *Config
}

type jwtClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
}

func NewJWTToken(config *Config) *JWTToken {
	return &JWTToken{
		config: config,
	}
}

func (j *JWTToken) CreateToken(userID int64) (string, error) {
	claims := jwtClaim{
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{},
		Exp:            time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.config.Signing_key))
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error) {
	log.Println("got token2 ", tokenString)
	log.Println("got token2 ", j.config.Signing_key)
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid auth token")
		}

		return []byte(j.config.Signing_key), nil
	})

	if err != nil {
		return 0, fmt.Errorf("invalid auth token")
	}

	claims, ok := token.Claims.(*jwtClaim)

	if !ok {
		return 0, fmt.Errorf("invalid auth token")
	}

	if claims.Exp < time.Now().Unix() {
		return 0, fmt.Errorf("auth token expired")
	}

	return claims.UserID, nil
}
