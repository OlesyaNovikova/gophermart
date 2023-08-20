package auth

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const tokenExp = time.Hour * 3

var key []byte

type claims struct {
	jwt.RegisteredClaims
	userID string
}

func InitAuth(k string) {
	key = []byte(k)
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(user string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		userID: user,
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserId проверяет валидность токена и возвращает UserId.
func GetUserId(tokenString string) (string, error) {
	cl := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, cl,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(key), nil
		})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		fmt.Println("Token is not valid")
		return "", err
	}
	fmt.Println("Token is valid")
	return cl.userID, nil
}
