update:
	@go mod tidy

impl_test:
	@make update
	@make codegen
	@go test ./...

server-test:
	@make codegen
	@cd test && go test

codegen:
	 @rm -rf ./internal/test/gen/api.gen.go
	 @oapi-codegen --config ./internal/test/resource/config.yaml ./internal/test/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./test/resource/local-config.yaml ./test/resource/pet-store.yaml  2>&1 | pbcopy

