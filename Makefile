test:
	@sh ./example/test

codegen:
	 @oapi-codegen --config ./config.yaml ./example/resource/pet-store.yaml

pet-store:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/pet-store.yaml  2>&1 | pbcopy

