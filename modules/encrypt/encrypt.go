package encrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword : this function will help to hash the controller input password to a
// secured encrypted string
func EncryptPassword(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while generating secured password")
		return "", err
	}

	return string(encryptPassword), nil
}

// VerifyEncryptPassword : this function will help to compare the encryted string with input password
// by the controller to during authentication
func VerifyEncryptPassword(password, hashedpassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}
	return true, nil
}
