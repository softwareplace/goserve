update:
	@go mod tidy

test:
	@make update
	@make codegen
	@go test ./...

codegen:
	 @rm -rf ./internal/gen/api.gen.go
	 @oapi-codegen --config ./internal/resource/config.yaml ./internal/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./internal/resource/config.yaml ./internal/resource/pet-store.yaml  2>&1 | xclip -selection clipboard

