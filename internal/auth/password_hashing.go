package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	password_hashed, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	return string(password_hashed), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
