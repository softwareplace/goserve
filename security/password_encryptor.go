package security

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

// Encrypt represents a utility for password hashing and validation.
type Encrypt struct {
	password        string
	encodedPassword string
	token           string
	authToken       string
	salt            string
}

// NewEncrypt creates a new Encrypt instance and generates hashed values.
func NewEncrypt(password string) *Encrypt {
	e := &Encrypt{password: password}
	e.encodedPassword = e.hashPassword(password)
	e.token = e.hashPassword(password + e.mixedString())
	e.authToken = e.hashPassword(password + e.mixedString())
	e.salt = e.hashPassword(password + e.mixedString())
	return e
}

// hashPassword hashes the given password using bcrypt.
func (e *Encrypt) hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), getBcryptCost())
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

// mixedString generates a random mixed string for additional entropy.
func (e *Encrypt) mixedString() string {
	randomResult, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	return fmt.Sprintf("%d", time.Now().UnixNano()+randomResult.Int64())
}

// IsValidPassword checks if the provided plaintext password matches the stored hash.
func (e *Encrypt) IsValidPassword(encodedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(e.password))
	return err == nil
}

// getBcryptCost retrieves the bcrypt cost from the environment variable B_CRYPT_COST.
// If the variable is not set or invalid, it defaults to bcrypt.DefaultCost.
func getBcryptCost() int {
	costStr := os.Getenv("B_CRYPT_COST")
	if costStr == "" {
		return bcrypt.DefaultCost
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil || cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return bcrypt.DefaultCost
	}
	return cost
}
