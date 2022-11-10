package token

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/yusuf/p-catalogue/pkg/encrypt"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    string
}

var signedKey = os.Getenv("secretKey")

//func GetSigningKey() string {
//	//
//	//keyString, err := encrypt.EncryptPassword(os.Getenv("secret_key"))
//	//if err != nil {
//	//	log.Fatalf("invalid secret key encryption")
//	//}
//	//return keyString
//}

func GenerateToken(id, email string) (string, string, error) {
	tokenClaims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "personal",
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		Email: email,
		ID:    id,
	}
	newtonClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		Issuer:    "personal",
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(signedKey))
	if err != nil {
		log.Println("cannot create token from claims")
		return "", "", err
	}
	newToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newtonClaims).SignedString([]byte(signedKey))
	if err != nil {
		log.Println("cannot create token from claims")
		return "", "", err
	}

	fmt.Println(token)
	return token, newToken, nil
}

func ParseTokenString(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		//if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		//	return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		//}
		return []byte(signedKey), nil
	})
	if err != nil {
		log.Fatalf("error while parsing token with it claims %v", err)
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		log.Fatalf("error %v user not authorized access", http.StatusUnauthorized)
	}
	if err := claims.Valid(); err != nil {
		log.Fatalf("error %v %s", http.StatusUnauthorized, err)
	}
	return claims, nil
}
