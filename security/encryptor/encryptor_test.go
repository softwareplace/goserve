package encryptor

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestNewEncrypt tests the NewEncrypt function.
func TestNewEncrypt(t *testing.T) {
	password := "mySecurePassword123"
	encrypt := NewEncrypt(password)

	if encrypt.EncodedPassword() == "" {
		t.Error("Expected encodedPassword to be set, got empty string")
	}

	if encrypt.Token() == "" {
		t.Error("Expected token to be set, got empty string")
	}

	if encrypt.Salt() == "" {
		t.Error("Expected salt to be set, got empty string")
	}
}

// TestHashPassword tests the hashPassword method.
func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"
	encrypt := &_PasswordEncryptorImpl{password: password}

	hashedPassword := encrypt.hashPassword(password)
	if hashedPassword == "" {
		t.Error("Expected hashedPassword to be set, got empty string")
	}

	// Verify that the hashed encryptor can be compared with the original encryptor
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Errorf("Failed to compare hashed encryptor: %v", err)
	}
}

// TestMixedString tests the mixedString method.
func TestMixedString(t *testing.T) {
	encrypt := &_PasswordEncryptorImpl{}

	mixedStr := encrypt.mixedString()
	if mixedStr == "" {
		t.Error("Expected mixedString to return a non-empty string, got empty string")
	}

	// Ensure the mixed string contains a timestamp and a random number
	_, err := strconv.ParseInt(mixedStr, 10, 64)
	if err != nil {
		t.Errorf("Failed to parse mixedString as int64: %v", err)
	}
}

// TestIsValidPassword tests the IsValidPassword method.
func TestIsValidPassword(t *testing.T) {
	password := "mySecurePassword123"
	encrypt := NewEncrypt(password)

	// Test with correct encryptor
	if !encrypt.IsValidPassword(encrypt.EncodedPassword()) {
		t.Error("Expected IsValidPassword to return true for correct encryptor")
	}

	// Test with incorrect encryptor
	incorrectPassword := "wrongPassword"
	if encrypt.IsValidPassword(incorrectPassword) {
		t.Error("Expected IsValidPassword to return false for incorrect encryptor")
	}
}

// TestGetBcryptCost tests the getBcryptCost function.
func TestGetBcryptCost(t *testing.T) {
	// Test default cost
	cost := getBcryptCost()
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected default cost to be %d, got %d", bcrypt.DefaultCost, cost)
	}

	// Test with valid custom cost
	os.Setenv("B_CRYPT_COST", "12")
	cost = getBcryptCost()
	if cost != 12 {
		t.Errorf("Expected cost to be 12, got %d", cost)
	}

	// Test with invalid custom cost (less than bcrypt.MinCost)
	os.Setenv("B_CRYPT_COST", "3")
	cost = getBcryptCost()
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected cost to be default (%d) for invalid value, got %d", bcrypt.DefaultCost, cost)
	}

	// Test with invalid custom cost (greater than bcrypt.MaxCost)
	os.Setenv("B_CRYPT_COST", "32")
	cost = getBcryptCost()
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected cost to be default (%d) for invalid value, got %d", bcrypt.DefaultCost, cost)
	}

	// Test with non-integer value
	os.Setenv("B_CRYPT_COST", "invalid")
	cost = getBcryptCost()
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected cost to be default (%d) for invalid value, got %d", bcrypt.DefaultCost, cost)
	}

	// Clean up environment variable
	os.Unsetenv("B_CRYPT_COST")
}

// TestMixedStringRandomness tests the randomness of the mixedString method.
func TestMixedStringRandomness(t *testing.T) {
	encrypt := &_PasswordEncryptorImpl{}

	// Generate two mixed strings and ensure they are different
	mixedStr1 := encrypt.mixedString()
	time.Sleep(1 * time.Millisecond) // Ensure timestamps are different
	mixedStr2 := encrypt.mixedString()

	if mixedStr1 == mixedStr2 {
		t.Error("Expected mixedString to generate different strings, got identical strings")
	}
}

// TestHashPasswordWithDifferentCosts tests the hashPassword method with different bcrypt costs.
func TestHashPasswordWithDifferentCosts(t *testing.T) {
	password := "mySecurePassword123"
	encrypt := &_PasswordEncryptorImpl{password: password}

	// Test with default cost
	hashedPassword1 := encrypt.hashPassword(password)

	// Test with custom cost
	os.Setenv("B_CRYPT_COST", "12")
	hashedPassword2 := encrypt.hashPassword(password)

	if hashedPassword1 == hashedPassword2 {
		t.Error("Expected different hashed passwords for different costs, got identical hashes")
	}

	// Clean up environment variable
	os.Unsetenv("B_CRYPT_COST")
}
