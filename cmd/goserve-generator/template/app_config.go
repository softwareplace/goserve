package template

const AppConfig = `package config

import "github.com/softwareplace/goserve/utils"

var (
	SwaggerFile = utils.GetEnvOrDefault(
		"SWAGGER_FILE",
		"./api/swagger.yaml",
	)

	ContextPath = utils.GetEnvOrDefault(
		"CONTEXT_PATH",
		"/api/${PROJECT}/v1/",
	)

	Port = utils.GetEnvOrDefault(
		"PORT",
		"8080",
	)
)
`
