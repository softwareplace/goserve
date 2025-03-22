package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/softwareplace/http-utils/api_context"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/server"
)

func main() {
	logger.LogSetup()

	server.Default().
		Get(func(ctx *api_context.ApiRequestContext[*api_context.DefaultContext]) {
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
		}, "/report/caller").
		StartServer()
}
