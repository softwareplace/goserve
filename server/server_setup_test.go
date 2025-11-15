package server

import (
	"os"
	"reflect"
	"time"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/login"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/internal/utils"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/secret"
)

func forTest(a Api[*goservectx.DefaultContext]) baseServer[*goservectx.DefaultContext] {
	baseServerValue := reflect.ValueOf(a).Elem()
	instance := baseServerValue.Interface()
	return instance.(baseServer[*goservectx.DefaultContext])
}

var (
	userPrincipalService *login.PrincipalServiceImpl
	securityService      security.Service[*goservectx.DefaultContext]

	loginService   *login.Service
	secretProvider *provider.TestSecretProviderImpl
	secretService  secret.Service[*goservectx.DefaultContext]
)

func beanInit() {
	userPrincipalService = login.NewPrincipalService()
	securityService = security.New(userPrincipalService)

	loginService = login.NewLoginService(securityService)
	secretProvider = provider.NewSecretProvider()

	secretService = secret.New(
		secretProvider,
		securityService,
	)
}

func testEnvSetup() {
	_ = os.Setenv("CONTEXT_PATH", "/")
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("API_SECRET_KEY", "DlJeR4%pPbB5Pr5cICMxg0xB")
	_ = os.Setenv("JWT_CLAIMS_ENCRYPTION_ENABLED", "true")
	_ = os.Setenv("API_PRIVATE_KEY", utils.TestSecretFilePath())
	beanInit()
}

func testEnvCleanup() {
	_ = os.Unsetenv("CONTEXT_PATH")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("API_SECRET_KEY")
	_ = os.Unsetenv("JWT_CLAIMS_ENCRYPTION_ENABLED")
	_ = os.Unsetenv("API_PRIVATE_KEY")
	_ = os.Unsetenv("FULL_AUTHORIZATION")
}

func getApiKey() (string, error) {
	response, err := securityService.From(provider.MockJWTSub, provider.MockScopes, time.Minute*10)
	if err != nil {
		return "", err
	}
	return response.JWT, nil
}
