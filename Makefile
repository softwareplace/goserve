test:
	@sh ./example/test

codegen:
	 @oapi-codegen --config ./resource/local-config.yaml ./example/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./example/resource/local-config.yaml ./example/resource/pet-store.yaml  2>&1 | pbcopy

