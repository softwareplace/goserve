package security

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"strings"
)

func HashPassword(password string) (string, string, error) {
	salt := GenerateSalt()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(hashedPassword), salt, nil
}

func GenerateSalt() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
