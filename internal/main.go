package main

import (
	"github.com/softwareplace/http-utils/internal/service/api"
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
		EmbeddedServer(api.Handler).
		Get(api.ReportCallerHandler, "/report/caller").
		StartServer()
}
