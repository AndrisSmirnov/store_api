package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(clientpassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(clientpassword), 16)
	return string(bytes), err
}

func CheckPasswordHash(passwordPost, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordPost))
	return err == nil
}
