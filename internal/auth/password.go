package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(password, digest string) error {
	err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(password))
	if err != nil {
		return err
	}

	return nil
}