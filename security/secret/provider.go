package secret

import apicontext "github.com/softwareplace/http-utils/context"

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
type Provider[T apicontext.Principal] interface {

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
	Get(ctx *apicontext.Request[T]) (string, error)
}
