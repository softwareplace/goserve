package login

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/service/provider"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/security/encryptor"
	"github.com/softwareplace/http-utils/security/jwt"
	"github.com/softwareplace/http-utils/security/login"
	"github.com/softwareplace/http-utils/security/principal"
	"github.com/softwareplace/http-utils/server"
	"sync"
	"time"
)

type PrincipalServiceImpl struct {
}

type Service struct {
	login.DefaultPasswordValidator[*apicontext.DefaultContext]
	securityService security.Service[*apicontext.DefaultContext]
}

var principalServiceService sync.Once
var onceLoginService sync.Once
var loginServiceInstance *Service
var principalServiceInstance *PrincipalServiceImpl

func NewLoginService(
	service security.Service[*apicontext.DefaultContext],
) *Service {
	onceLoginService.Do(func() {
		loginServiceInstance = &Service{
			securityService: service,
		}
	})
	return loginServiceInstance
}

func NewPrincipalService() principal.Service[*apicontext.DefaultContext] {
	principalServiceService.Do(func() {
		principalServiceInstance = &PrincipalServiceImpl{}
	})
	return principalServiceInstance
}

func (d *PrincipalServiceImpl) LoadPrincipal(ctx *apicontext.Request[*apicontext.DefaultContext]) bool {
	if ctx.Authorization == "" {
		return false

	}

	context := apicontext.NewDefaultCtx()
	context.SetRoles("api:key:generator")
	ctx.Principal = &context
	return true
}

func (l *Service) RequiredScopes() []string {
	return []string{
		"api:key:generator",
	}
}

func (l *Service) GetApiJWTInfo(apiKeyEntryData server.ApiKeyEntryData,
	_ *apicontext.Request[*apicontext.DefaultContext],
) (jwt.Entry, error) {
	return jwt.Entry{
		Client:     apiKeyEntryData.ClientName,
		Key:        apiKeyEntryData.ClientId,
		Expiration: apiKeyEntryData.Expiration,
		Scopes: []string{
			"api:example:user",
			"api:example:admin",
			"read:pets",
			"write:pets",
			"api:key:generator",
		},
	}, nil
}

func (l *Service) OnGenerated(data jwt.Response,
	jwtEntry jwt.Entry,
	ctx apicontext.SampleContext[*apicontext.DefaultContext],
) {
	provider.MockStore[jwtEntry.Key] = *jwtEntry.PublicKey
	log.Printf("%s - %s", jwtEntry.Key, data.Token)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (l *Service) Login(user login.User) (*apicontext.DefaultContext, error) {
	result := &apicontext.DefaultContext{}
	result.SetRoles("api:example:user", "api:example:admin", "read:pets", "write:pets", "api:key:generator")
	password := encryptor.NewEncrypt(user.Password).EncodedPassword()
	result.SetEncryptedPassword(password)
	return result, nil
}

func (l *Service) TokenDuration() time.Duration {
	return time.Minute * 15
}

func (l *Service) SecurityService() security.Service[*apicontext.DefaultContext] {
	return l.securityService
}
