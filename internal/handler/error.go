package handler

import (
	log "github.com/sirupsen/logrus"
	goservecontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/server"
	"sync"
)

type errorHandlerImpl struct {
}

var (
	errorHandlerInstance goservecontext.ApiHandler[*goservecontext.DefaultContext]
	errorHandlerOnce     sync.Once
)

func New() goservecontext.ApiHandler[*goservecontext.DefaultContext] {
	errorHandlerOnce.Do(func() {
		errorHandlerInstance = &errorHandlerImpl{}
	})
	return errorHandlerInstance
}

func (p *errorHandlerImpl) Handler(ctx *goservecontext.Request[*goservecontext.DefaultContext], err error, source string) {
	log.Errorf("%s failed with error: %+v", source, err)
	if source == server.ErrorHandlerWrapper {
		ctx.InternalServerError("Internal server error")
	}

	if source == server.SecurityValidatorResourceAccess {
		ctx.Unauthorized()
	}
}
