package encryptor

type Service interface {

	// Secret retrieves the secret key used to sign and validate JWT tokens.
	// This function ensures consistent access to the secret key across the pService.
	//
	// Returns:
	//   - A byte slice containing the secret key.
	Secret() []byte

	// Encrypt encrypts the given value using the secret associated with the apiSecurityServiceImpl instance.
	// It returns the encrypted string or an error if encryption fails.
	Encrypt(key string) (string, error)

	// Decrypt decrypts the given encrypted string using the secret associated with the apiSecurityServiceImpl instance.
	// It returns the decrypted string or an error if decryption fails.
	//
	// Parameters:
	// - encrypted: The string that has been encrypted and needs to be decrypted.
	//
	// Returns:
	// - A string representing the decrypted value if the operation is successful.
	// - An error if decryption fails due to issues like invalid cipher text or incorrect secret.
	//
	// Notes:
	// - The decryption logic must use secure cryptographic mechanisms to ensure data safety.
	// - Ensure that any sensitive data involved in the decryption process is handled securely
	//   and not exposed in logs or error messages.
	Decrypt(encrypted string) (string, error)
}
