package login

import (
	log "github.com/sirupsen/logrus"
	goservecontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/internal/service/provider"
	"github.com/softwareplace/goserve/security"
	"github.com/softwareplace/goserve/security/encryptor"
	"github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/login"
	"github.com/softwareplace/goserve/security/principal"
	"github.com/softwareplace/goserve/server"
	"sync"
	"time"
)

type PrincipalServiceImpl struct {
}

type Service struct {
	login.DefaultPasswordValidator[*goservecontext.DefaultContext]
	securityService security.Service[*goservecontext.DefaultContext]
}

var principalServiceService sync.Once
var onceLoginService sync.Once
var loginServiceInstance *Service
var principalServiceInstance *PrincipalServiceImpl

func NewLoginService(
	service security.Service[*goservecontext.DefaultContext],
) *Service {
	onceLoginService.Do(func() {
		loginServiceInstance = &Service{
			securityService: service,
		}
	})
	return loginServiceInstance
}

func NewPrincipalService() principal.Service[*goservecontext.DefaultContext] {
	principalServiceService.Do(func() {
		principalServiceInstance = &PrincipalServiceImpl{}
	})
	return principalServiceInstance
}

func (d *PrincipalServiceImpl) LoadPrincipal(ctx *goservecontext.Request[*goservecontext.DefaultContext]) bool {
	if ctx.Authorization == "" {
		return false
	}

	context := goservecontext.NewDefaultCtx()
	//context.SetRoles("api:key:generator")
	context.SetRoles("api:key:generator", "write:pets", "read:pets")
	ctx.Principal = &context
	return true
}

func (l *Service) RequiredScopes() []string {
	return []string{
		"api:key:generator",
	}
}

func (l *Service) GetApiJWTInfo(apiKeyEntryData server.ApiKeyEntryData,
	_ *goservecontext.Request[*goservecontext.DefaultContext],
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
	ctx goservecontext.SampleContext[*goservecontext.DefaultContext],
) {
	provider.MockStore[jwtEntry.Key] = *jwtEntry.PublicKey
	log.Printf("%s - %s", jwtEntry.Key, data.Token)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (l *Service) Login(user login.User) (*goservecontext.DefaultContext, error) {
	result := &goservecontext.DefaultContext{}
	//result.SetRoles("api:example:user", "api:example:admin", "read:pets", "write:pets", "api:key:generator")
	result.SetRoles("api:example:user")
	password := encryptor.NewEncrypt(user.Password).EncodedPassword()
	result.SetEncryptedPassword(password)
	return result, nil
}

func (l *Service) TokenDuration() time.Duration {
	return time.Minute * 15
}

func (l *Service) SecurityService() security.Service[*goservecontext.DefaultContext] {
	return l.securityService
}
