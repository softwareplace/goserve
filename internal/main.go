package main

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/http-utils/context"
	"github.com/softwareplace/http-utils/logger"
	"github.com/softwareplace/http-utils/server"
)

func init() {
	// Setup log system. Using nested-logrus-formatter -> https://github.com/antonfisher/nested-logrus-formatter?tab=readme-ov-file
	// Reload log file target reference based on `LOG_FILE_NAME_DATE_FORMAT`
	logger.LogSetup()
}

func main() {
	server.Default().
		Get(ReportCallerHandler, "/report/caller").
		StartServer()
}

func ReportCallerHandler(ctx *apicontext.ApiRequestContext[*apicontext.DefaultContext]) {
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
