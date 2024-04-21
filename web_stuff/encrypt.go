package httpstuff

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
)

func EncryptPassword(password, key string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"pwd": password,
	})
	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		log.Printf("err occured: %v", err)
		panic(err)
	}
	return tokenStr
}
