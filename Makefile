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

