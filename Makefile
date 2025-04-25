update:
	@go mod tidy

test:
	@make update
	@make codegen
	@go test -v -bench=. ./...

codegen:
	 @rm -rf ./internal/gen/api.gen.go
	 @oapi-codegen --config ./internal/resource/config.yaml ./internal/resource/pet-store.yaml

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
