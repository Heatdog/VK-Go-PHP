package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenFileds struct {
	ID string `json:"sub"`
}

func GenerateToken(fields TokenFileds, key string) (string, error) {
	payload := jwt.MapClaims{
		"sub": fields.ID,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(key))
}

func VerifyToken(tokenString, key string) (*TokenFileds, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token calims are not of type *TokenClaims")
	}

	res := &TokenFileds{
		ID: claims["sub"].(string),
	}

	return res, nil
}
