package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/go-password/pkg/str"

	"github.com/softwareplace/goserve/env"
	"github.com/softwareplace/goserve/internal/service/apiservice"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/logger"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/secret"
	"github.com/softwareplace/goserve/server"
)

func init() {
	// Setup log system. Using nested-logrus-formatter -> https://github.com/antonfisher/nested-logrus-formatter?tab=readme-ov-file
	// Reload log file target reference based on `LOG_FILE_NAME_DATE_FORMAT`
	logger.LogSetup()

	if secretKey := env.GetEnvOrDefault("API_SECRET_KEY", ""); secretKey == "" {
		randomString := str.New().
			Generate()

		log.Infof("API_SECRET_KEY: %s", randomString)
		_ = os.Setenv("API_SECRET_KEY", randomString)
		_ = os.Setenv("API_PRIVATE_KEY", "./internal/resource/secret/private.key")
	}

}

func runSecretApi() {
	userPrincipalService := login.NewPrincipalService()
	securityService := security.New(
		userPrincipalService,
	)

	loginService := login.NewLoginService(securityService)
	secretProvider := provider.NewSecretProvider()

	secretService := secret.New(
		secretProvider,
		securityService,
	).DisableForPublicPath(true)

	server.Default().
		LoginResourceEnabled(true).
		SecretKeyGeneratorResourceEnabled(true).
		LoginService(loginService).
		SecretService(secretService).
		SecurityService(securityService).
		EmbeddedServer(apiservice.Register).
		Get(apiservice.ReportCallerHandler, "/report/caller").
		SwaggerDocHandler("./internal/resource/pet-store.yaml").
		StartServer()
}

func runPublicApi() {
	userPrincipalService := login.NewPrincipalService()
	securityService := security.New(
		userPrincipalService,
	)

	loginService := login.NewLoginService(securityService)

	server.Default().
		LoginService(loginService).
		SecurityService(securityService).
		SwaggerDocHandler("./internal/resource/pet-store.yaml").
		EmbeddedServer(apiservice.Register).
		Get(apiservice.ReportCallerHandler, "/report/caller").
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
