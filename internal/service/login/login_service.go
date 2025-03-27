package login

import (
	log "github.com/sirupsen/logrus"
	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/login"
	"github.com/softwareplace/goserve/security/principal"
	"github.com/softwareplace/goserve/security/secret"
	"sync"
	"time"
)

type PrincipalServiceImpl struct {
}

type Service struct {
	login.DefaultPasswordValidator[*goservectx.DefaultContext]
	securityService security.Service[*goservectx.DefaultContext]
}

var principalServiceService sync.Once
var onceLoginService sync.Once
var loginServiceInstance *Service
var principalServiceInstance *PrincipalServiceImpl

func NewLoginService(
	service security.Service[*goservectx.DefaultContext],
) *Service {
	onceLoginService.Do(func() {
		loginServiceInstance = &Service{
			securityService: service,
		}
	})
	return loginServiceInstance
}

func NewPrincipalService() principal.Service[*goservectx.DefaultContext] {
	principalServiceService.Do(func() {
		principalServiceInstance = &PrincipalServiceImpl{}
	})
	return principalServiceInstance
}

func (d *PrincipalServiceImpl) LoadPrincipal(ctx *goservectx.Request[*goservectx.DefaultContext]) bool {
	if ctx.Authorization == "" {
		return false
	}

	context := goservectx.NewDefaultCtx()
	//context.SetRoles("api:key:generator")
	context.SetRoles("api:key:generator", "write:pets", "read:pets")
	ctx.Principal = &context
	return true
}

func (l *Service) RequiredScopes() []string {
	return []string{
		"api:key:generator",
		"api:key:generator-v2",
	}
}

func (l *Service) GetApiJWTInfo(apiKeyEntryData secret.ApiKeyEntryData,
	_ *goservectx.Request[*goservectx.DefaultContext],
) (secret.Entry, error) {
	return secret.Entry{
		Key:        apiKeyEntryData.ClientId,
		Expiration: apiKeyEntryData.Expiration,
		Roles: []string{
			"api:example:user",
			"api:example:admin",
			"read:pets",
			"write:pets",
			"api:key:generator",
		},
	}, nil
}

func (l *Service) OnGenerated(data jwt.Response,
	jwtEntry secret.Entry,
	ctx goservectx.SampleContext[*goservectx.DefaultContext],
) {
	provider.MockStore[jwtEntry.Key] = *jwtEntry.PublicKey
	log.Printf("%s - %s", jwtEntry.Key, data.JWT)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (l *Service) Login(user login.User) (*goservectx.DefaultContext, error) {
	result := &goservectx.DefaultContext{}
	//result.SetRoles("api:example:user", "api:example:admin", "read:pets", "write:pets", "api:key:generator")
	result.SetRoles("api:example:user")
	password := encryptor.NewEncrypt(user.Password).EncodedPassword()
	result.SetEncryptedPassword(password)
	return result, nil
}

func (l *Service) TokenDuration() time.Duration {
	return time.Minute * 15
}

func (l *Service) SecurityService() security.Service[*goservectx.DefaultContext] {
	return l.securityService
}
