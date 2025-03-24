package handler

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	errorhandler "github.com/softwareplace/goserve/error"
	"github.com/softwareplace/goserve/server"
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

func (p *errorHandlerImpl) Handler(ctx *apicontext.Request[*apicontext.DefaultContext], err error, source string) {
	log.Errorf("%s failed with error: %+v", source, err)
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
