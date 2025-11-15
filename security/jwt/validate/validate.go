package validate

type Validate interface {

	// IsValid checks if the provided JWT token string is valid.
	// It parses the token string using the configured secret key and verifies the token's validity.
	//
	// Parameters:
	//   - tokenString: The JWT token string to be validated.
	//
	// Returns:
	//   - True if the token is successfully parsed and is valid; otherwise, false.
	IsValid(tokenString string) bool
}
