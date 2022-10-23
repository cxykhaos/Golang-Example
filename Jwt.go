package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("khaos")

type User struct {
	User_Id String
}

type Claims struct {
	User_Id string
	jwt.StandardClaims
}

func ReleaseToken(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	claims := &Claims{
		User_Id: user.User_Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "khaos",
			Subject:   "user token",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
