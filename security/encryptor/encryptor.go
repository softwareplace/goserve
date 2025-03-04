package encryptor

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

// PasswordEncryptor is an interface for securely hashing and validating passwords using bcrypt.
// It provides methods to generate hashed passwords, tokens, and salts, as well as to validate passwords.
//
// Environment Variables:
//   - B_CRYPT_COST The cost factor for bcrypt hashing.  Default 10
type PasswordEncryptor interface {
	// EncodedPassword returns the hashed version of the password.
	EncodedPassword() string

	// Token returns a hashed token generated from the password and additional entropy.
	Token() string

	// Salt returns a hashed salt generated from the password and additional entropy.
	Salt() string

	// IsValidPassword checks if the provided plaintext password matches the stored hash.
	IsValidPassword(encodedPassword string) bool
}

type _PasswordEncryptorImpl struct {
	password        string
	encodedPassword string
	token           string
	salt            string
}

func (e *_PasswordEncryptorImpl) EncodedPassword() string {
	return e.encodedPassword
}

func (e *_PasswordEncryptorImpl) Token() string {
	return e.token
}

func (e *_PasswordEncryptorImpl) Salt() string {
	return e.salt
}

// NewEncrypt creates a new PasswordEncryptor instance and generates hashed values for the password, token, and salt.
//
// Parameters:
//   - password: The plaintext password to be hashed.
//
// Returns:
//   - A PasswordEncryptor instance with the hashed password, token, and salt.
func NewEncrypt(password string) PasswordEncryptor {
	e := &_PasswordEncryptorImpl{password: password}
	e.encodedPassword = e.hashPassword(password)
	e.token = e.hashPassword(password + e.mixedString())
	e.salt = e.hashPassword(password + e.mixedString())
	return e
}

// hashPassword hashes the given password using bcrypt.
//
// Parameters:
//   - password: The plaintext password to be hashed.
//
// Returns:
//   - A string containing the hashed password.
//
// Notes:
//   - The bcrypt cost is determined by the B_CRYPT_COST environment variable or defaults to bcrypt.DefaultCost.
func (e *_PasswordEncryptorImpl) hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), getBcryptCost())
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

// mixedString generates a random mixed string for additional entropy.
//
// Returns:
// 	- A string containing a timestamp and a cryptographically secure random number.
//
// Notes:
// 	- This method is used to add randomness to the token and salt generation process.

func (e *_PasswordEncryptorImpl) mixedString() string {
	randomResult, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}
	return fmt.Sprintf("%d", time.Now().UnixNano()+randomResult.Int64())
}

// IsValidPassword checks if the provided plaintext password matches the stored hash.
//
// Parameters:
// - encodedPassword: The hashed password to compare against.
//
// Returns:
//   - A boolean indicating whether the plaintext password matches the hashed password.
//
// Notes:
//   - This method uses bcrypt.CompareHashAndPassword for secure comparison.
func (e *_PasswordEncryptorImpl) IsValidPassword(encodedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(e.password))
	return err == nil
}

// getBcryptCost retrieves the bcrypt cost from the environment variable B_CRYPT_COST.
// If the variable is not set or invalid, it defaults to bcrypt.DefaultCost.
//
// Returns:
// - An integer representing the bcrypt cost.
//
// Notes:
// - The cost must be between bcrypt.MinCost and bcrypt.MaxCost.
// - If the environment variable is not set or invalid, the default cost is used.

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
