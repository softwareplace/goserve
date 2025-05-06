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
	@go install github.com/softwareplace/goserve/cmd/goserve-generator@v0.0.1-SNAPSHOT

goserve-generator:
	@go run cmd/goserve-generator/main.go -n go-example -u softwareplace -r true
	@cd go-example && git status
	@rm -rf go-example

install-coverage-utils:
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/matm/gocov-html/cmd/gocov-html@latest
