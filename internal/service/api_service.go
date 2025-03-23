package service

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/internal/gen"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/server"
)

func baseResponse(message string, status int) gen.BaseResponse {
	success := false
	timestamp := 1625867200

	response := gen.BaseResponse{
		Message:   &message,
		Code:      &status,
		Success:   &success,
		Timestamp: &timestamp,
	}
	return response
}

type ApiService struct {
	UserService
	PetService
	FileService
	InventoryService
}

func ApiServiceHandler(handler server.Api[*apicontext.DefaultContext]) {
	handler.EmbeddedServer(gen.ApiResourceHandler(&ApiService{})).
		SwaggerDocHandler("./internal/resource/pet-store.yaml")
}

func ReportCallerHandler(ctx *apicontext.Request[*apicontext.DefaultContext]) {
	enable := ctx.QueryOf("enable")
	if enable == "true" {
		logger.LogReportCaller = true
		log.SetReportCaller(true)
		ctx.Ok(map[string]interface{}{
			"message": "Logger report caller enabled",
		})
	} else {
		logger.LogReportCaller = false
		log.SetReportCaller(false)
		ctx.Ok(map[string]interface{}{
			"message": "Logger report caller disabled",
		})
	}
}
