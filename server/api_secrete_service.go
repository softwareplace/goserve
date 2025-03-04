package server

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/error_handler"
	"github.com/softwareplace/http-utils/security"
	"log"
	"os"
	"time"
)

type ApiKeyEntryData struct {
	ClientName string        `json:"clientName"` // Client information for which the public key is generated (required)
	Expiration time.Duration `json:"expiration"` // Expiration specifies the duration until the API key expires (optional).
	ClientId   string        `json:"clientId"`   // ClientId represents the unique identifier for the client associated with the API key entry (required).
}

type ApiKeyGeneratorService[T api_context.ApiPrincipalContext] interface {

	// SecurityService returns an instance of ApiSecurityService responsible for handling security-related operations.
	// This includes operations such as JWT generation, claims extraction, encryption, decryption, and authorization handling.
	// It provides the foundational security mechanisms required by the ApiKeyGeneratorService.
	//
	// Returns:
	//   - security.ApiSecurityService[T]: The security service instance associated with the implementing service,
	//	 providing security functionalities for API keys, JWTs, and authorization processes.
	SecurityService() security.ApiSecurityService[T]

	// GetApiJWTInfo generates the security.ApiJWTInfo for the given ApiKeyEntryData and ApiRequestContext.
	// This method is responsible for processing the API key entry data and request context to create an ApiJWTInfo object,
	// which contains essential JWT-related information such as the client, key, and expiration details.
	//
	// Parameters:
	//   - apiKeyEntryData: An instance of ApiKeyEntryData that includes client details, expiration duration, and unique client identifier.
	//   - ctx: The API request context, which contains metadata and principal information related to the API key generation process.
	//
	// Returns:
	//   - security.ApiJWTInfo: The generated ApiJWTInfo object containing JWT details necessary for creating the API secret JWT.
	//   - error: If an error occurs during the process, it returns the corresponding error; otherwise, nil.
	GetApiJWTInfo(apiKeyEntryData ApiKeyEntryData, ctx *api_context.ApiRequestContext[T]) (security.ApiJWTInfo, error)

	// OnGenerated is invoked after an API key has been successfully generated.
	// This function allows additional processing or handling, such as logging,
	// auditing, or notifying dependent systems of the newly generated API key.
	//
	// Parameters:
	//   - data: The generated token as security.JwtResponse.
	//   - apiJWTInfo: The requested security.ApiJWTInfo.
	//   - ctx: The API request context, containing metadata and principal
	//		  information related to the API key generation.
	OnGenerated(data security.JwtResponse, apiJWTInfo security.ApiJWTInfo, ctx api_context.SampleContext[T])

	// RequiredScopes defines the list of scopes that are required to access the API key generator functionality.
	//
	// This method returns a slice of strings that represent the necessary authorization
	// scopes required for clients to access this service. These scopes ensure fine-grained
	// access control and enforce security policies.
	//
	// Returns:
	//   - []string: A slice of strings representing the required scopes.
	//	 For example, scope identifiers such as "read:apikey", "write:apikey" can
	//	 be included here to indicate necessary permissions.
	RequiredScopes() []string
}

func (a *apiRouterHandlerImpl[T]) apiKeyGeneratorDataHandler(ctx *api_context.ApiRequestContext[T], apiKeyEntryData ApiKeyEntryData) {
	error_handler.Handler(func() {
		log.Printf("API/KEY/GENERATOR: requested by: %s", ctx.AccessId)

		jwtInfo, err := a.apiKeyGeneratorService.GetApiJWTInfo(apiKeyEntryData, ctx)

		if err != nil {
			log.Printf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		if jwtInfo.PublicKey == nil || *jwtInfo.PublicKey == "" {
			key, err := a.generatePubKey(a.apiSecretAccessHandler.SecretKey())
			if err != nil {
				log.Printf("API/KEY/GENERATOR: Failed to generate public key: %v", err)
				ctx.InternalServerError("Failed to generate JWT. Please try again later.")
				return
			}

			jwtInfo.PublicKey = &key
		}

		jwtInfo.Expiration = time.Hour * jwtInfo.Expiration
		jwt, err := a.apiSecurityService.GenerateApiSecretJWT(jwtInfo)

		if err != nil {
			log.Printf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		ctx.Ok(jwt)

		a.apiKeyGeneratorService.OnGenerated(*jwt, jwtInfo, ctx.GetSample())
	}, func(err error) {
		log.Printf("API/KEY/GENERATOR/HANDLER: Failed to handle request: %v", err)
		ctx.InternalServerError("Failed to generate JWT. Please try again later.")
	})
}

// generatePubKey generates an encrypted public key from a given private key file.
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
func (a *apiRouterHandlerImpl[T]) generatePubKey(secretKey string) (string, error) {
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

	encryptedKey, err := a.apiSecurityService.Encrypt(string(publicKeyPEM))

	if err != nil {
		log.Fatalf("Failed to encryptor public key: %s", err)
		return "", nil
	}
	return encryptedKey, err
}
