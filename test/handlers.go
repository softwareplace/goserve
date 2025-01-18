package main

import (
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/server"
)

type principalServiceImpl struct {
}

func (p *principalServiceImpl) Roles(ctx api_context.ApiRequestContext[*api_context.DefaultContext]) []string {
	return []string{"admin"}
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
