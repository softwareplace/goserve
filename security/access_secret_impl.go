package security

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security/principal"
	"log"
	"net/http"
	"os"
)

type apiSecretHandlerImpl[T api_context.ApiPrincipalContext] struct {
	service                        ApiSecurityService[T]
	provider                       ApiSecretKeyProvider[T]
	principalService               principal.PService[T]
	secretKey                      string
	apiSecret                      any
	ignoreValidationForPublicPaths bool
}

// ApiSecretAccessHandlerBuild apiSecretHandlerImpl is an implementation of the ApiSecretAccessHandler interface which manages
// security-related operations for API requests, such as validating API keys and initializing
// cryptographic keys. It encapsulates the logic for validating an API secret key and restricting
// unauthorized access to resources.
//
// Type Parameters:
//   - T: A type that satisfies the `api_context.ApiPrincipalContext` interface, providing API principal-specific context.
//
// Fields:
//   - secretKey: The file path to the secret key used for cryptographic operations.
//   - service: An instance of ApiSecurityService responsible for cryptographic and security services.
//   - apiSecret: Holder of the parsed private key, supporting either ECDSA or RSA key types.
//   - provider: A provider responsible for loading the API secret key for access validation.
//   - principalService: A service managing API principal claims and IDs to ensure request security.
//
// This struct provides methods to initialize the secret key, validate the public key against the private key,
// and enforce access security middleware, ensuring requests are authorized with proper credentials.
func ApiSecretAccessHandlerBuild[T api_context.ApiPrincipalContext](
	secretKey string,
	provider ApiSecretKeyProvider[T],
	service ApiSecurityService[T],
) ApiSecretAccessHandler[T] {

	handler := apiSecretHandlerImpl[T]{
		secretKey: secretKey,
		service:   service,
		provider:  provider,
	}
	handler.initAPISecretKey()
	return &handler
}

func (a *apiSecretHandlerImpl[T]) SecretKey() string {
	return a.secretKey
}

func (a *apiSecretHandlerImpl[T]) DisableForPublicPath(ignore bool) ApiSecretAccessHandler[T] {
	a.ignoreValidationForPublicPaths = ignore
	return a
}

func (a *apiSecretHandlerImpl[T]) HandlerSecretAccess(ctx *api_context.ApiRequestContext[T]) bool {
	if a.ignoreValidationForPublicPaths && principal.IsPublicPath[T](*ctx) {
		return true
	}

	if !a.apiSecretKeyValidation(ctx) {
		a.service.handlerErrorOrElse(ctx, nil, ApiSecretAccessHandlerError, func() {
			ctx.Error("You are not allowed to access this resource", http.StatusUnauthorized)
		})
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
// - Stores the private key in the `apiSecret` field of the `apiSecretHandlerImpl` struct.
//
// Logs an error and terminates the application if any of the above steps fail.
func (a *apiSecretHandlerImpl[T]) initAPISecretKey() {
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
// - Loads the API secret using the provided `ApiSecretKeyServiceProvider`.
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
func (a *apiSecretHandlerImpl[T]) apiSecretKeyValidation(ctx *api_context.ApiRequestContext[T]) bool {
	// Decode the Base64-encoded public key
	claims, err := a.service.JWTClaims(ctx)

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %v", err)
		return false
	}

	ctx.ApiKeyClaims = claims

	apiKey, err := a.service.Decrypt(claims["apiKey"].(string))

	if err != nil {
		log.Printf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %v", err)
		return false
	}

	ctx.ApiKeyId = apiKey

	apiAccessKey, err := a.provider.Get(ctx)
	if err != nil {
		log.Printf("API_SECRET_LOADER: AuthorizationHandler failed: %v", err)
		return false
	}

	// Decode the PEM-encoded public key
	decryptKey, err := a.service.Decrypt(apiAccessKey)
	if err != nil {
		log.Printf("API_SECRET_DECRYPT: AuthorizationHandler failed: %v", err)
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

// GeneratePubKey generates an encrypted public key from a given private key file.
//
// This function performs the following steps:
// - Reads the private key from the specified file path.
// - Decodes the PEM block from the private key data.
// - Parses the private key using the PKCS8 format.
// - Determines the type of the private key (ECDSA or RSA).
// - Marshals the corresponding public key into PEM format.
// - Encrypts the generated PEM-encoded public key using the service's encryption logic.
//
// Arguments:
//   - secretKey (string): The file path to the private key.
//
// Returns:
//   - (string, error): An encrypted PEM-encoded public key and an error (if any occurred).
//
// Errors:
//   - Fails if the private key file cannot be read, parsed, or if the key type is unsupported.
//   - Fails if the public key cannot be marshaled or encrypted.
//
// Example:
//
//	 encryptedPubKey, err := handler.generatePubKey("path/to/private.key")
//	 if err != nil {
//		 log.Printf("Error generating public key: %v", err)
//	 }
func (a *apiSecretHandlerImpl[T]) GeneratePubKey(secretKey string) (string, error) {
	privateKeyData, err := os.ReadFile(secretKey)
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

	// Generate and log the corresponding public key
	var publicKeyBytes []byte
	switch key := privateKey.(type) {
	case *ecdsa.PrivateKey:
		log.Println("Loaded ECDSA private key successfully")
		publicKeyBytes, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			log.Fatalf("Failed to marshal ECDSA public key: %s", err.Error())
		}
	case *rsa.PrivateKey:
		log.Println("Loaded RSA private key successfully")
		publicKeyBytes, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			log.Fatalf("Failed to marshal RSA public key: %s", err.Error())
		}
	default:
		log.Fatalf("Unsupported private key type: %T", key)
	}

	// Encode the public key in PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	encryptedKey, err := a.service.Encrypt(string(publicKeyPEM))

	if err != nil {
		log.Fatalf("Failed to encrypt public key: %s", err)
		return "", nil
	}
	return encryptedKey, err
}
