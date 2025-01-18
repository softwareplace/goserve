package impl

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/security"
	"github.com/softwareplace/http-utils/server"
	"time"
)

type loginServiceImpl struct {
	securityService security.ApiSecurityService[*api_context.DefaultContext]
}

func (l *loginServiceImpl) SecurityService() security.ApiSecurityService[*api_context.DefaultContext] {
	return l.securityService
}

func New(securityService security.ApiSecurityService[*api_context.DefaultContext]) server.LoginService[*api_context.DefaultContext] {
	return &loginServiceImpl{
		securityService: securityService,
	}
}

func (l *loginServiceImpl) Login(user server.LoginEntryData) (*api_context.DefaultContext, error) {
	result := &api_context.DefaultContext{}
	result.SetRoles("api:example:user", "api:example:admin")
	return result, nil
}

func (l *loginServiceImpl) TokenDuration() time.Duration {
	return time.Minute * 15
}
