package auth

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/server"
	"log"
	"net/http"
	"os"
)

type ApiSecretKeyLoader[T api_context.ApiContextData] func(ctx *api_context.ApiRequestContext[T]) (string, error)

type ApiSecretAccessHandler[T api_context.ApiContextData] interface {
}

type apiSecurityHandlerImpl[T api_context.ApiContextData] struct {
	secretKey string
	service   *security.ApiSecurityService[T]
	apiSecret any
	loader    ApiSecretKeyLoader[T]
}

// Handler creates a new instance of ApiSecretAccessHandler, which is a middleware security handler
// responsible for managing API security using a private key and an optional loader for resolving API secret keys.
//
// Args:
//
//	  secretKey (string):
//		 - A file path to the private key file that will be used to initialize the API secret.
//		 - The private key must be in PKCS8 format and supported types are ECDSA and RSA.
//	  loader (ApiSecretKeyLoader[T]):
//		 - A function responsible for loading the API secret key, which is often used to decrypt or verify API-related data.
//		 - Can be nil if no custom secret key loading functionality is needed.
//	  service (*security.ApiSecurityService[T]):
//		 - An instance of ApiSecurityService used for JWT claim extraction, encryption, and decryption operations.
//		 - This service acts as the main utility for security tasks within the API.
//
// Returns:
//
//	  ApiSecretAccessHandler[T]:
//		 - An implementation of the ApiSecretAccessHandler interface that provides security validation middleware.
//
// Example Usage:
//
//	service := &security.ApiSecurityService[YourContext]{}
//	handler := Handler("path/to/private-key.pem", yourSecretKeyLoader, service)
//	isAuthorized := handler.apiSecretMiddleware(apiRequestContext)
func Handler[T api_context.ApiContextData](
	secretKey string,
	loader ApiSecretKeyLoader[T],
	service *security.ApiSecurityService[T],
	routerHandler server.ApiRouterHandler[T],
) {
	handler := &apiSecurityHandlerImpl[T]{
		secretKey: secretKey,
		service:   service,
		loader:    loader,
	}
	handler.initAPISecretKey()
	routerHandler.Use(handler.apiSecretMiddleware, "API_SECRET_MIDDLEWARE")
}

func (a *apiSecurityHandlerImpl[T]) apiSecretMiddleware(ctx *api_context.ApiRequestContext[T]) bool {
	if !a.apiSecretKeyValidation(ctx) {
		ctx.Error("You are not allowed to access this resource", http.StatusUnauthorized)
		return false
	}
	return true
}

// initAPISecretKey initializes the API secret key by reading and parsing a private key file.
//
// This function performs the following steps:
// - Reads the private key file specified by the `secretKey` field.
// - Decodes the PEM block from the file data.
// - Parses the private key from the PEM data using PKCS8 format.
// - Validates the type of the private key (either ECDSA or RSA).
// - Stores the private key in the `apiSecret` field of the `apiSecurityHandlerImpl` struct.
//
// Logs an error and terminates the application if any of the above steps fail.
func (a *apiSecurityHandlerImpl[T]) initAPISecretKey() {
	// Load private key from the provided secretKey file path
	privateKeyData, err := os.ReadFile(a.secretKey)
	if err != nil {
		log.Fatalf("Failed to read private key file: %s", err.Error())
	}

	// Decode PEM block from the private key data
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatalf("Failed to decode private key PEM block")
	}

	// Parse the private key using ParsePKCS8PrivateKey
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err.Error())
	}
	a.apiSecret = privateKey

	switch key := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		log.Println("Loaded ECDSA private key successfully")
	case *rsa.PrivateKey:
		log.Println("Loaded RSA private key successfully")
	default:
		log.Fatalf("Unsupported private key type: %T", key)
	}
}

// apiSecretKeyValidation verifies the validity of a public key against the private key stored in the handler.
//
// This function performs the following steps:
// - Extracts JWT claims from the request context using the ApiSecurityService.
// - Loads the API secret using the provided `ApiSecretKeyLoader`.
// - Decrypts the API access key to retrieve the PEM-encoded public key.
// - Decodes the PEM-encoded public key and parses it into a usable public key object.
// - Validates the type of the parsed public key (ECDSA or RSA).
// - Ensures the extracted public key corresponds to the private key stored in the `apiSecret` field.
//
// If any of the above steps fail, the function logs the error and returns `false`, indicating that
// the public key validation has failed. Otherwise, it returns `true`.
//
// Args:
//
//	  ctx (*api_context.ApiRequestContext[T]):
//		 - The context of the API request carrying the necessary data for validation.
//
// Returns:
//
//	  bool:
//		 - `true` if the public key is valid and corresponds to the private key.
//		 - `false` if the public key is invalid or the validation fails.
func (a *apiSecurityHandlerImpl[T]) apiSecretKeyValidation(ctx *api_context.ApiRequestContext[T]) bool {
	// Decode the Base64-encoded public key
	claims, err := (*a.service).JWTClaims(*ctx)

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: Authorization failed: %v", err)
		return false
	}
	apiContext := ctx.RequestData

	apiContext.SetApiKeyClaims(claims)

	apiKey, err := (*a.service).Decrypt(claims["apiKey"].(string))

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: Authorization failed: %v", err)
		return false
	}

	apiContext.SetApiKeyId(apiKey)

	apiAccessKey, err := a.loader(ctx)
	if err != nil {
		log.Printf("API_SECRET_LOADER: Authorization failed: %v", err)
		return false
	}

	// Decode the PEM-encoded public key
	decryptKey, err := (*a.service).Decrypt(apiAccessKey)
	if err != nil {
		log.Printf("API_SECRET_DECRYPT: Authorization failed: %v", err)
		return false
	}
	block, _ := pem.Decode([]byte(decryptKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Printf("Failed to decode public key PEM block")
		return false
	}

	// Parse the public key
	parsedPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Printf("Failed to parse public key: %v", err)
		return false
	}

	switch privateKey := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		// Ensure the type of the public key matches ECDSA
		publicKey, ok := parsedPublicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Printf("invalid public key type, expected ECDSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.X.Cmp(privateKeyPubKey.X) != 0 || publicKey.Y.Cmp(privateKeyPubKey.Y) != 0 {
			log.Printf("public key does not match the private key")
			return false
		}
	case *rsa.PrivateKey:
		// Ensure the type of the public key matches RSA
		publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
		if !ok {
			log.Printf("invalid public key type, expected RSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.E != privateKeyPubKey.E || publicKey.N.Cmp(privateKeyPubKey.N) != 0 {
			log.Printf("public key does not match the private key")
			return false
		}
	default:
		log.Printf("unsupported private key type: %T", privateKey)
		return false
	}

	return true
}
