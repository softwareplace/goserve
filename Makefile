test:
	@sh ./example/test

codegen:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/swagger.yaml

codegen-x-api-key:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/swagger-with-api-key.yaml


