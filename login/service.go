package login

import (
	apicontext "github.com/softwareplace/http-utils/context"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Service[T apicontext.Principal] interface {

	// Login processes the login request for the specified user by validating their credentials.
	// It authenticates the user based on the provided login data and returns an authenticated principal context or an error.
	//
	// Parameters:
	//   - user: An instance of User that contains the username, encryptor, and/or email for user authentication.
	//
	// Returns:
	//   - T: The authenticated principal context representing the logged-in user.
	//   - error: If authentication fails, an error is returned.
	Login(user User) (T, error)

	// TokenDuration specifies the duration for which a generated JWT token remains valid.
	// This value defines the time-to-live (TTL) for the token, ensuring secure and proper session management.
	//
	// Returns:
	//   - time.Duration: The duration for which a generated token is valid.
	TokenDuration() time.Duration

	// IsValidPassword validates the user-provided plaintext password against the stored encrypted password.
	//
	// This method uses the encryptor package to create a password hash from the provided loginEntryData password
	// and compares it with the encrypted password available in the principal context.
	//
	// A default implementation is available as server.DefaultPasswordValidator[T]
	//
	// Parameters:
	//   - loginEntryData: The User containing the plaintext password to be validated.
	//   - principal: The principal context of type T, which contains the stored encrypted password.
	//
	// Returns:
	//   - bool: True if the passwords match; false otherwise.
	IsValidPassword(loginEntryData User, principal T) bool
}
