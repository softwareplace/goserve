package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

type principalServiceImpl struct {
}

func (d *principalServiceImpl) LoadPrincipal(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) (*api_context.DefaultContext, bool) {
	return &api_context.DefaultContext{}, true
}

func (d *principalServiceImpl) SetData(*api_context.DefaultContext) {

}

func (d *principalServiceImpl) GetRoles(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) []string {
	return []string{"admin", "GET", "test:v2"}
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

func (p *errorHandlerImpl) Handler(ctx *api_context.ApiRequestContext[*api_context.DefaultContext], err error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
