package template

const Makefile = `update:
	@go mod tidy

test:
	@make update
	@make codegen
	@go test -v -bench=. ./...

codegen:
	 @mkdir -p ./internal/adapter/handler/gen/ || true
	 @rm -rf ./internal/adapter/handler/gen/** || true
	 @touch ./internal/adapter/handler/gen/api.gen.go
	 @oapi-codegen --config ./config/config.yaml ./api/swagger.yaml

`
