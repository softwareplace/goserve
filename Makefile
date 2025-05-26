APP_NAME=go-example
update:
	@./setup

coverage:
	@./coverage-validator

run-test:
	@make update
	@go test -v -bench=. ./...
	@make goserve-generator

pet-store:
	 @oapi-codegen --config ./internal/resource/config.yaml ./internal/resource/pet-store.yaml  2>&1 | xclip -selection clipboard

# Try test implementation
run:
	@make update
	@make codegen
	@PROTECTED_API=true LOG_REPORT_CALLER=true go run ./internal/main.go

install:
	@go install github.com/softwareplace/goserve/cmd/goserve-generator@$(shell git rev-parse HEAD)

goserve-generator:
	@go run cmd/goserve-generator/main.go -n $(APP_NAME) -u softwareplace -r true -cgf internal/resource/test_config.yaml -gsv $(shell git rev-parse HEAD)
	@cd $(APP_NAME) && git status


install-coverage-utils:
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/matm/gocov-html/cmd/gocov-html@latest
