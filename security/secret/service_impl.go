package secret

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	goserveerror "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/encryptor"
	goservejwt "github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/model"
	"github.com/softwareplace/goserve/security/router"
	"github.com/softwareplace/goserve/utils"
	"net/http"
	"os"
	"time"
)

// New creates and initializes a new instance of Service interface implementation,
// which handles security-related operations for API requests, such as validating API keys and initializing
// cryptographic keys. It encapsulates the logic necessary for validating API secret keys, restricting unauthorized
// resource access, and managing cryptographic operations.
//
// Type Parameters:
//   - T: A type that satisfies the `goservectx.Principal` interface, representing the API request's security context.
//
// Parameters:
//   - provider: An implementation of the Provider interface responsible for managing API secrets and key-related operations.
//   - service: An implementation of the security.Service interface for cryptographic and validation-related services.
//
// Returns:
//   - An instance of the Service interface, initialized and ready to perform security operations.
//
// Notes:
//   - The function retrieves the API private key from the "API_PRIVATE_KEY" environment variable.
//   - It terminates the application if the private key is not set or fails to initialize properly.
func New[T goservectx.Principal](
	provider Provider[T],
	service security.Service[T],
) Service[T] {
	secretKey := utils.GetEnvOrDefault("API_PRIVATE_KEY", "")

	if secretKey == "" {
		log.Panicf("API_PRIVATE_KEY environment variable not set")
	}

	handler := apiSecretHandlerImpl[T]{
		secretKey: secretKey,
		Service:   service,
		Provider:  provider,
	}
	handler.InitAPISecretKey()
	return &handler
}

func (a *apiSecretHandlerImpl[T]) Handler(ctx *goservectx.Request[T], apiKeyEntryData model.ApiKeyEntryData) {
	goserveerror.Handler(func() {
		log.Infof("API/KEY/GENERATOR: requested by: %s", ctx.AccessId)

		info, err := a.GetJwtEntry(apiKeyEntryData, ctx)

		if err != nil {
			log.Errorf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		if info.PublicKey == nil || *info.PublicKey == "" {
			var key string

			key, err = a.GeneratePubKey(a.SecretKey())
			if err != nil {
				log.Errorf("API/KEY/GENERATOR: Failed to generate public key: %v", err)
				ctx.InternalServerError("Failed to generate JWT. Please try again later.")
				return
			}

			info.PublicKey = &key
		}

		info.Expiration = time.Hour * info.Expiration

		response, err := a.From(info.Key, info.Roles, info.Expiration)

		if err != nil {
			log.Errorf("API/KEY/GENERATOR: Failed to generate JWT: %v", err)
			ctx.InternalServerError("Failed to generate JWT. Please try again later.")
			return
		}

		ctx.Ok(response)

		a.OnGenerated(*response, info, ctx.GetSample())
	}, func(err error) {
		log.Errorf("API/KEY/GENERATOR/HANDLER: Failed to handle request: %v", err)
		ctx.InternalServerError("Failed to generate JWT. Please try again later.")
	})
}

func (a *apiSecretHandlerImpl[T]) SecretKey() string {
	return a.secretKey
}

func (a *apiSecretHandlerImpl[T]) DisableForPublicPath(ignore bool) Service[T] {
	a.ignoreValidationForPublicPaths = ignore
	return a
}

func (a *apiSecretHandlerImpl[T]) HandlerSecretAccess(ctx *goservectx.Request[T]) bool {
	isPublicPath := router.IsPublicPath[T](*ctx)
	if a.ignoreValidationForPublicPaths && isPublicPath {
		return true
	}

	if !a.ApiSecretKeyValidation(ctx) {
		a.HandlerErrorOrElse(ctx, nil, AccessHandlerError, nil)
		ctx.Error("You are not allowed to access this resource", http.StatusUnauthorized)
		return false
	}
	return true
}

func (a *apiSecretHandlerImpl[T]) InitAPISecretKey() {
	// Load private key from the provided secretKey file path
	privateKeyData, err := os.ReadFile(a.secretKey)
	if err != nil {
		log.Panicf("Failed to read private key file: %+v", err)
	}

	// Decode PEM block from the private key data
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Panicf("Failed to decode private key PEM block")
	}

	// Parse the private key using ParsePKCS8PrivateKey
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Panicf("Failed to parse private key: %+v", err)
	}
	a.apiSecret = privateKey

	switch key := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		log.Println("Loaded ECDSA private key successfully")
	case *rsa.PrivateKey:
		log.Println("Loaded RSA private key successfully")
	default:
		log.Panicf("Unsupported private key type: %T", key)
	}
}

func (a *apiSecretHandlerImpl[T]) ApiSecretKeyValidation(ctx *goservectx.Request[T]) bool {
	// Decode the Base64-encoded public key
	claims, err := a.Decode(ctx.ApiKey)

	if err != nil {
		log.Errorf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %+v", err)
		return false
	}

	ctx.ApiKeyClaims = claims

	isJwtClaimsEncryptionEnabled := encryptor.JwtClaimsEncryptionEnabled()

	apiKey := claims[goservejwt.SUB].(string)

	if isJwtClaimsEncryptionEnabled {
		apiKey, err = a.Decrypt(claims[goservejwt.SUB].(string))
	}

	if err != nil {
		log.Errorf("JWT/CLAIMS_EXTRACT: AuthorizationHandler failed: %v", err)
		return false
	}

	ctx.ApiKeyId = apiKey

	apiAccessKey, err := a.GetPublicKey(ctx)
	if err != nil {
		log.Errorf("API_SECRET_LOADER: AuthorizationHandler failed: %v", err)
		return false
	}

	// Decode the PEM-encoded public key
	decryptKey, err := a.Decrypt(apiAccessKey)
	if err != nil {
		log.Errorf("API_SECRET_DECRYPT: AuthorizationHandler failed: %v", err)
		return false
	}

	block, _ := pem.Decode([]byte(decryptKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Errorf("Failed to decode public key PEM block")
		return false
	}

	// Parse the public key
	parsedPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorf("Failed to parse public key: %v", err)
		return false
	}

	switch privateKey := a.apiSecret.(type) {
	case *ecdsa.PrivateKey:
		// Ensure the type of the public key matches ECDSA
		publicKey, ok := parsedPublicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Errorf("invalid public key type, expected ECDSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.X.Cmp(privateKeyPubKey.X) != 0 || publicKey.Y.Cmp(privateKeyPubKey.Y) != 0 {
			log.Errorf("public key does not match the private key")
			return false
		}
	case *rsa.PrivateKey:
		// Ensure the type of the public key matches RSA
		publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
		if !ok {
			log.Errorf("invalid public key type, expected RSA")
			return false
		}

		// Validate if the public key corresponds to the private key
		privateKeyPubKey := &privateKey.PublicKey
		if publicKey.E != privateKeyPubKey.E || publicKey.N.Cmp(privateKeyPubKey.N) != 0 {
			log.Errorf("public key does not match the private key")
			return false
		}
	default:
		log.Errorf("unsupported private key type: %T", privateKey)
		return false
	}

	return true
}

func (a *apiSecretHandlerImpl[T]) GeneratePubKey(secretKey string) (string, error) {
	privateKeyData, err := os.ReadFile(secretKey)
	if err != nil {
		log.Panicf("Failed to read private key file: %+v", err)
	}

	// Decode PEM block from the private key data
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Panicf("Failed to decode private key PEM block")
	}

	// Parse the private key using ParsePKCS8PrivateKey
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Panicf("Failed to parse private key: %+v", err)
	}

	// Generate and log the corresponding public key
	var publicKeyBytes []byte
	switch key := privateKey.(type) {
	case *ecdsa.PrivateKey:
		log.Println("Loaded ECDSA private key successfully")
		publicKeyBytes, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			log.Panicf("Failed to marshal ECDSA public key: %+v", err)
		}
	case *rsa.PrivateKey:
		log.Println("Loaded RSA private key successfully")
		publicKeyBytes, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			log.Panicf("Failed to marshal RSA public key: %+v", err)
		}
	default:
		log.Panicf("Unsupported private key type: %T", key)
	}

	// Encode the public key in PEM format
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return a.Encrypt(string(publicKeyPEM))
}
