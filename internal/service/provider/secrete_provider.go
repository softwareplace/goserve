package provider

import (
	log "github.com/sirupsen/logrus"

	goservectx "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/jwt/response"
	"github.com/softwareplace/goserve/security/model"
)

type TestSecretProviderImpl struct {
	TestProviderGet func(ctx *goservectx.Request[*goservectx.DefaultContext]) (string, error)
	TestGetJwtEntry func(apiKeyEntryData model.ApiKeyEntryData,
		_ *goservectx.Request[*goservectx.DefaultContext],
	) (model.Entry, error)

	TestOnGenerated func(data response.Response,
		jwtEntry model.Entry,
		ctx goservectx.SampleContext[*goservectx.DefaultContext],
	)

	TestRequiredScopes func() []string
}

var MockScopes = []string{
	"api:example:user",
	"api:example:admin",
	"read:pets",
	"write:pets",
	"api:key:goserve-generator",
}

var MockJWTSub = "37c75552614a4eb58a2eb2d04928cdfd"
var MockStore = map[string]string{
	MockJWTSub: "SDEBLUQofvUky3K9q9EaBHFZLU2Xvizng4uYRaEBszR7MQ4hQa1CRnOcM6xh/3+sneyiRSFiMn4EJX08j0w8tu3x06wEjzRpY68izbqT2F9ToPGrGkmrDEplvEPuIqlrvi/l7MyjJ+4T/Elhue2Kzqfjo6TF6n6Vuju1wVDC8Y5hLNmW2/5vbWPrkYaeAysnza9jq52CRC9JeJC+TFke0AoyWUKeu7wRoH9ygu0RsK2ZrBJ/K2llXHIZ55FIv+D58+NoztxFGOgNBw+DcvAdVCykD2TgDA2wxXGUbng4Fmv0pAst12waCYNvpblWfFELkukAZ0xNxnaeX6sO/dwLL0qUgvDppHs+W1goS7UUmIN9tzf/vNtUTgda5BYK8sDIH9lmOqWQq59aaaXO1/TTbpZ2xpHNk2IB25G6Z0WQIxD+26KbzrQ07aSr5cV1ff7wHMIGAoLEKqJqCkYud2Z52Ss/v5/9fNzWwlQ/IWtHzFXAjWmf+8I3Olh7X74kvhzUN4miA3evXS5PI94hOubtDXhd6w4SOYN6CXuZ+RQYllESaziBWUf5jBo487CwMgKQXwvzuEY1oulFUn3zjptzMXa2L/g6UY2zV2CZbnnhB+lzprg=",
}

func NewSecretProvider() *TestSecretProviderImpl {
	return &TestSecretProviderImpl{}
}

func (s *TestSecretProviderImpl) GetPublicKey(ctx *goservectx.Request[*goservectx.DefaultContext]) (string, error) {
	if s.TestProviderGet != nil {
		return s.TestProviderGet(ctx)
	}

	return MockStore[ctx.ApiKeyId], nil
}

func (s *TestSecretProviderImpl) GetJwtEntry(apiKeyEntryData model.ApiKeyEntryData,
	ctx *goservectx.Request[*goservectx.DefaultContext],
) (model.Entry, error) {

	if s.TestGetJwtEntry != nil {
		return s.TestGetJwtEntry(apiKeyEntryData, ctx)
	}

	return model.Entry{
		Key:        apiKeyEntryData.ClientId,
		Expiration: apiKeyEntryData.Expiration,
		Roles:      MockScopes,
	}, nil
}

func (s *TestSecretProviderImpl) OnGenerated(data response.Response,
	jwtEntry model.Entry,
	ctx goservectx.SampleContext[*goservectx.DefaultContext],
) {
	if s.TestOnGenerated != nil {
		s.TestOnGenerated(data, jwtEntry, ctx)
		return
	}

	MockStore[jwtEntry.Key] = *jwtEntry.PublicKey
	log.Printf("%s - %s", jwtEntry.Key, data.JWT)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (s *TestSecretProviderImpl) RequiredScopes() []string {
	if s.TestRequiredScopes != nil {
		return s.TestRequiredScopes()
	}
	return []string{
		"api:example:user",
		"api:example:admin",
		"api:key:goserve-generator",
	}
}
