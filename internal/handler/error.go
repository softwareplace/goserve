package handler

import (
	apicontext "github.com/softwareplace/http-utils/context"
	errorhandler "github.com/softwareplace/http-utils/error"
	"github.com/softwareplace/http-utils/server"
	"sync"
)

type errorHandlerImpl struct {
}

var (
	errorHandlerInstance errorhandler.ApiHandler[*apicontext.DefaultContext]
	errorHandlerOnce     sync.Once
)

func New() errorhandler.ApiHandler[*apicontext.DefaultContext] {
	errorHandlerOnce.Do(func() {
		errorHandlerInstance = &errorHandlerImpl{}
	})
	return errorHandlerInstance
}

func (p *errorHandlerImpl) Handler(ctx *apicontext.Request[*apicontext.DefaultContext], _ error, source string) {
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
