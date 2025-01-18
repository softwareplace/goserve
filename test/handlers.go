package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

type principalServiceImpl struct {
}

func (d *principalServiceImpl) LoadPrincipal(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) bool {
	ctx.Principal = &api_context.DefaultContext{}
	return true
}

func (d *principalServiceImpl) SetAuthorizationClaims(map[string]interface{}) {

}

func (d *principalServiceImpl) SetApiKeyClaims(map[string]interface{}) {

}

func (d *principalServiceImpl) SetApiKeyId(string) {

}

func (d *principalServiceImpl) SetAccessId(string) {

}

type errorHandlerImpl struct {
}

func (p *errorHandlerImpl) Handler(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], _ error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
