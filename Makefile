test:
	@sh ./example/test

codegen:
	 @oapi-codegen --config ./resource/oapi-codegen-config.yaml ./example/resource/swagger.yaml


