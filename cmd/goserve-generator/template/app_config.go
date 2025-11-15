package template

const AppConfig = `package config

import goserverenv "github.com/softwareplace/goserve/env"

var (
	SwaggerFile = goserverenv.GetEnvOrDefault(
		"SWAGGER_FILE",
		"./api/swagger.yaml",
	)

	ContextPath = goserverenv.GetEnvOrDefault(
		"CONTEXT_PATH",
		"/api/${PROJECT}/v1/",
	)

	Port = goserverenv.GetEnvOrDefault(
		"PORT",
		"8080",
	)
)
`
