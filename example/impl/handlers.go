package impl

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

type PrincipalServiceImpl struct {
}

func (d *PrincipalServiceImpl) LoadPrincipal(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) bool {
	ctx.Principal = &api_context.DefaultContext{}
	return true
}

func (d *PrincipalServiceImpl) SetAuthorizationClaims(map[string]interface{}) {

}

func (d *PrincipalServiceImpl) SetApiKeyClaims(map[string]interface{}) {

}

func (d *PrincipalServiceImpl) SetApiKeyId(string) {

}

func (d *PrincipalServiceImpl) SetAccessId(string) {

}

type ErrorHandlerImpl struct {
}

func (p *ErrorHandlerImpl) Handler(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], _ error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
