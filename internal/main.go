package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/http-utils/internal/handler"
	"github.com/softwareplace/http-utils/internal/service/api"
	"github.com/softwareplace/http-utils/internal/service/login"
	"github.com/softwareplace/http-utils/internal/service/provider"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/secret"
	"github.com/softwareplace/http-utils/server"
	"os"
)

func init() {
	// Setup log system. Using nested-logrus-formatter -> https://github.com/antonfisher/nested-logrus-formatter?tab=readme-ov-file
	// Reload log file target reference based on `LOG_FILE_NAME_DATE_FORMAT`
	logger.LogSetup()
}

func runSecretApi() {
	userPrincipalService := login.NewPrincipalService()
	errorHandler := handler.New()
	securityService := security.New(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
		errorHandler,
	)

	loginService := login.NewLoginService(securityService)
	secretProvider := provider.NewSecretProvider()

	secretHandler := secret.New(
		"./internal/secret/private.key",
		secretProvider,
		securityService,
	)

	server.Default().
		LoginResourceEnabled(true).
		SecretKeyGeneratorResourceEnabled(true).
		ApiKeyGeneratorResource(loginService).
		LoginService(loginService).
		SecretService(secretHandler).
		SecurityService(securityService).
		PrincipalService(userPrincipalService).
		EmbeddedServer(api.Handler).
		SwaggerDocHandler("./internal/resource/pet-store.yaml").
		Get(api.ReportCallerHandler, "/report/caller").
		StartServer()
}

func runPublicApi() {
	userPrincipalService := login.NewPrincipalService()
	errorHandler := handler.New()
	securityService := security.New(
		"ue1pUOtCGaYS7Z1DLJ80nFtZ",
		userPrincipalService,
		errorHandler,
	)

	loginService := login.NewLoginService(securityService)

	server.Default().
		LoginResourceEnabled(true).
		LoginService(loginService).
		SecurityService(securityService).
		PrincipalService(userPrincipalService).
		EmbeddedServer(api.Handler).
		SwaggerDocHandler("./internal/resource/pet-store.yaml").
		Get(api.ReportCallerHandler, "/report/caller").
		StartServer()
}

func main() {
	if env, _ := os.LookupEnv("PROTECTED_API"); env == "true" {
		log.Info("Protected API enabled")
		runSecretApi()
	} else {
		log.Warning("Protected API disabled.")
		runPublicApi()
	}
}
