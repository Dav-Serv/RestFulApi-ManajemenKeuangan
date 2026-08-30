package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID	uint	`json:"user_id"`
	Email	string	`json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT yang berlaku selama expiryHours jam.
func GenerateToken(userID uint, email string, secret string, expiryHours string) (string, error) {
	hours, err := strconv.Atoi(expiryHours)
	if err != nil || hours <= 0 {
		hours = 72
	}

	Claims := Claims{
		UserID: userID,
		Email: 	email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt:	jwt.NewNumericDate(time.Now().Add(time.Duration(hours) * time.Hour)),
			IssuedAt:	jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken memverifikasi signature & masa berlaku token, lalu mengembalikan claims-nya.
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("method signing token tidak valid")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}