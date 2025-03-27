package secret

import (
	goservectx "github.com/softwareplace/goserve/context"
	goservejwt "github.com/softwareplace/goserve/security/jwt"
	"time"
)

type ApiKeyEntryData struct {
	ClientName string        `json:"clientName"` // Client information for which the public key is generated (required)
	Expiration time.Duration `json:"expiration"` // Expiration specifies the duration until the API key expires (optional).
	ClientId   string        `json:"clientId"`   // ClientId represents the unique identifier for the client associated with the API key entry (required).
}

type Entry struct {
	Key        string
	Expiration time.Duration
	Roles      []string
	PublicKey  *string
}

// Provider is an interface designed to provide secure access to API secret keys
// based on the context of an incoming API request.
//
// This interface plays a crucial role in enabling secure communication and access control
// by retrieving API keys that are specific to each request. It abstracts the mechanism for
// obtaining these keys, which could involve fetching from a database, a configuration file,
// or any other secure storage mechanism.
//
// Type Parameters:
//   - T: A type that satisfies the `context.Principal` interface, representing
//     the authentication and authorization context for API requests.
type Provider[T goservectx.Principal] interface {

	// Get (ctx *context.Request[T]) (string, error):
	//	   Fetches the API secret key for the given request context. The method should implement
	//	   any necessary logic to securely retrieve and provide the key, such as decryption or
	//	   validation.
	//
	// Example Use Case:
	// When processing an API request that requires validation with a secret key, the implementation
	// of this interface can retrieve and provide the appropriate key tailored to the request's context.
	//
	// Returns:
	//   - A string representing the API secret key.
	//   - An error if the key retrieval or processing fails, ensuring proper error handling in the
	//	 request lifecycle.
	Get(ctx *goservectx.Request[T]) (string, error)

	// GetJwtEntry generates the jwt.Entry for the given ApiKeyEntryData and Request.
	// This method is responsible for processing the API key entry data and request context to create an ApiJWTInfo object,
	// which contains essential JWT-related information such as the client, key, and expiration details.
	//
	// Parameters:
	//   - apiKeyEntryData: An instance of ApiKeyEntryData that includes client details, expiration duration, and unique client identifier.
	//   - ctx: The API request context, which contains metadata and principal information related to the API key generation process.
	//
	// Returns:
	//   - jwt.Entry: The generated ApiJWTInfo object containing JWT details necessary for creating the API secret JWT.
	//   - error: If an error occurs during the process, it returns the corresponding error; otherwise, nil.
	GetJwtEntry(apiKeyEntryData ApiKeyEntryData, ctx *goservectx.Request[T]) (Entry, error)

	// OnGenerated is invoked after an API key has been successfully generated.
	// This function allows additional processing or handling, such as logging,
	// auditing, or notifying dependent systems of the newly generated API key.
	//
	// Parameters:
	//   - response: The generated token as jwt.Response.
	//   - jwtEntry: The requested jwt.Entry.
	//   - ctx: The API request context, containing metadata and principal
	//		  information related to the API key generation.
	OnGenerated(response goservejwt.Response, jwtEntry Entry, ctx goservectx.SampleContext[T])

	// RequiredScopes returns a list of scopes that are mandatory for accessing the related end point
	// resources registered by invoking server.Api SecretService method.
	RequiredScopes() []string
}
