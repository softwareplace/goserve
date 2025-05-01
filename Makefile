update:
	@make codegen
	@go mod tidy

coverage:
	@mkdir .out || true
	@go test ./... -coverprofile=.out/coverage_raw.out
	@grep -v "/test/" .out/coverage_raw.out > .out/coverage.out
	@go tool cover -func=.out/coverage.out
	@go tool cover -html=.out/coverage.out -o .out/index.html

run-test:
	@make update
	@go test -v -bench=. ./...
	@make goserve-generator

codegen:
	 @rm -rf ./test/gen/api.gen.go || true
	 @mkdir -p ./test/gen || true
	 @oapi-codegen --config ./test/resource/config.yaml ./test/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./test/resource/config.yaml ./test/resource/pet-store.yaml  2>&1 | xclip -selection clipboard

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

