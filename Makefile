test:
	@sh ./example/test

codegen:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/swagger.yaml

pet-store:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/pet-store.yaml  2>&1 | pbcopy


codegen-x-api-key:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/swagger-with-api-key.yaml


