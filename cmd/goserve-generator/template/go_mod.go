package template

const GoMod = `module github.com/${USERNAME}/${PROJECT}

go 1.24.2

require (
	github.com/google/uuid v1.6.0
	github.com/oapi-codegen/runtime v1.1.1
	github.com/sirupsen/logrus v1.9.3
	github.com/softwareplace/goserve v1.0.4
)
`
