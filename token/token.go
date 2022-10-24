package token

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yusuf/p-catalogue/pkg/encrypt"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Email string
	ID    string
}

func GetSigningKey() []byte {
	keyString, err := encrypt.EncryptPassword(os.Getenv("secret_key"))
	if err != nil {
		log.Fatalf("invalid secret key encryption")
	}
	return []byte(keyString)
}

func GenerateToken(id, email string) (interface{}, interface{}, error) {
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
	newtokenClaims := &jwt.RegisteredClaims{
		Issuer:    "personal",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, tokenClaims).SignedString(GetSigningKey())
	if err != nil {
		log.Println("cannot create token from claims")
		return "", "", err
	}
	newToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, newtokenClaims).SignedString(GetSigningKey())
	if err != nil {
		log.Println("cannot create token from claims")
		return "", "", err
	}

	return token, newToken, nil
}

func ParseTokenString(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		}
		return GetSigningKey(), nil
	})
	if err != nil {
		log.Fatalf("error while parsing token with it claims")
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
