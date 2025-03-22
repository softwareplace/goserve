package context

// DefaultRequestId is a default requester ID used when no custom requester ID is provided.
// This is a hardcoded value and is **not secure for production use**.
// In production, you should use a unique and secure requester ID for each request.
var DefaultRequestId = "081162586c7f4f77b877fbca0f09cb7f"

// DefaultContext is a simple implementation of a context that holds roles, an encrypted password, and an optional requester ID.
// This struct is intended for use in **test and development environments only** and is **not recommended for production use**.
//
// In production environments, you should use a more secure and robust implementation that:
//   - Properly manages sensitive data (e.g., passwords).
//   - Validates and sanitizes input.
//   - Integrates with your authentication and authorization systems.
//
// Fields:
//   - roles: A slice of strings representing the roles associated with this context.
//   - encryptedPassword: A string representing the encrypted password. Note that this is not secure for production use.
//   - DefaultRequesterId: A pointer to a string representing an optional requester ID. If not provided, the default requester ID is used.
type DefaultContext struct {
	roles              []string
	encryptedPassword  string
	DefaultRequesterId *string
}

// NewDefaultCtx creates and returns a new instance of DefaultContext.
// This is a convenience function for initializing the context in test/development environments.
func NewDefaultCtx() *DefaultContext {
	return &DefaultContext{}
}

// RequesterId returns the requester ID associated with the context.
// If no custom requester ID is set, it returns the default requester ID.
// WARNING: The default requester ID is hardcoded and is **not secure for production use**.
// In production, you should use a unique and secure requester ID for each request.
func (d *DefaultContext) RequesterId() string {
	if d.DefaultRequesterId == nil || *d.DefaultRequesterId == "" {
		return DefaultRequestId
	}
	return *d.DefaultRequesterId
}

// GetRoles returns the roles associated with the context.
// This is useful for role-based access control (RBAC).
func (d *DefaultContext) GetRoles() []string {
	return d.roles
}

// SetEncryptedPassword sets the encrypted password in the context.
// WARNING: This method does not perform any encryption or validation.
// In production, you should use a secure encryption mechanism (e.g., bcrypt, Argon2).
func (d *DefaultContext) SetEncryptedPassword(encryptedPassword string) {
	d.encryptedPassword = encryptedPassword
}

// EncryptedPassword returns the encrypted password stored in the context.
// WARNING: This method does not decrypt the password. It simply returns the stored value.
func (d *DefaultContext) EncryptedPassword() string {
	return d.encryptedPassword
}

// SetRoles sets the roles associated with the context.
// If no roles are provided, it initializes the roles slice as empty.
// This is useful for simulating role-based access control in test/development environments.
func (d *DefaultContext) SetRoles(roles ...string) {
	if roles == nil {
		d.roles = []string{}
	}
	d.roles = append(d.roles, roles...)
}
