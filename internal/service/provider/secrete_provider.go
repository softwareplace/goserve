package provider

import (
	log "github.com/sirupsen/logrus"
	apicontext "github.com/softwareplace/goserve/context"
	"github.com/softwareplace/goserve/security/jwt"
	"github.com/softwareplace/goserve/security/secret"
	"sync"
)

type secretProviderImpl []struct{}

var MockStore = map[string]string{
	"37c75552614a4eb58a2eb2d04928cdfd": "Kt96yFYCuIjqgJlDtufWpdoUdwNNRBEsqoWqjZn456h/dzxdKhAuG3F98pU6d2PaIQ4H16tgN0Rw7nwzVV6YsyL2y03ms9TlydbXNx/fqrVJoN9U9MVvpwV6CuorHhAv+6C9Wguw9LUByqVzX36IuV/agOrWuHW9RrHchwV1OktMtFYWPuoLWRDID20RYqtaYSPx1VYaftfSolgPjG3a5oKHOtxuTt9fz8guVroZdFbaLF8B8qigvp3/YXT/+xmaPIeXJrIObb0i4Kkw2w24SzI8W3c5bzJjqWKYec4OGVibBOE+ZwL0UpP2mib2MOZ2pRauChRVdlV/4shFEUwOpnFvyT8xwEgIkgvMKKl7xPbUEZTfRif4I7jzJ/NS3Xfx0qd4jorlOyAgY5uDq+Ijlte89OfvZg27O0l9iQ+ezwwhzavbHU8xhVnq5MsbPZHQGvM+p3aELe1RQg43e/JohLvASXM1rmA17mx7eNbB6BAkuaz1E8BcGLHvNQpyO+T4zGY9agqJ6stF7aatNNShQQ5tN31bbKHyynH0BZtoE3E9j8pXiWsx2Xs8hd0lWA3wguslbZAN65JlSeT4PbRKDiFfTl+vQRUQoJu+aAY+OpOnduM=",
	/// X-Api-Key for test only eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiRWZUZElvbm8vd3o2ZUxOWG5PajZOQ2lOUnpqTWNMSklYdWlEeUt5d1ZEST0iLCI1elA0cWs1a2Q2aUtCekZ2Y2NiTUI3OW5EbUEwczgrY0dsMHVZT2s4MUE5cCIsIjdMVnZDVTlXbVl2SVY2OU1sTHdIZHpXb0hlV0VSSlBpQ1E9PSIsInNNem8vYjlUTGVHMVBwUjFkYkV5MGhmRC9vbHZkalZpeVIwPSIsIngzODhLdTkxdUJHTncwckp1MHcyRVhIR0JZajVKVUVaZFBuV2g0b1JyMk1rIl0sImV4cCI6NTE5OTA5ODAxMywiaWF0IjoxNzQzMDk4MDEzLCJpc3MiOiJnb3NlcnZlci1leGFtcGxlIiwic3ViIjoibFdQKzdHTjNzZjhoNVZXcVRyaTBUM0RaSHNaYmEvWWcwenV4TWhKK0o4Mkw2R0FHelRkUFl6N2hGV0doWkhBYiJ9.6-Z4W5np8uXLuQJttd9BOvuG7iG9EFC8RsTL2fB0OqU
}

func (s *secretProviderImpl) Get(ctx *apicontext.Request[*apicontext.DefaultContext]) (string, error) {
	return MockStore[ctx.ApiKeyId], nil
}

var secreteProvideOnce sync.Once
var secretProvider secret.Provider[*apicontext.DefaultContext]

func NewSecretProvider() secret.Provider[*apicontext.DefaultContext] {
	secreteProvideOnce.Do(func() {
		secretProvider = &secretProviderImpl{}
	})
	return secretProvider
}

func (s *secretProviderImpl) GetJwtEntry(apiKeyEntryData secret.ApiKeyEntryData,
	_ *apicontext.Request[*apicontext.DefaultContext],
) (secret.Entry, error) {
	return secret.Entry{
		Key:        apiKeyEntryData.ClientId,
		Expiration: apiKeyEntryData.Expiration,
		Roles: []string{
			"api:example:user",
			"api:example:admin",
			"read:pets",
			"write:pets",
			"api:key:generator",
		},
	}, nil
}

func (s *secretProviderImpl) OnGenerated(data jwt.Response,
	jwtEntry secret.Entry,
	ctx apicontext.SampleContext[*apicontext.DefaultContext],
) {
	MockStore[jwtEntry.Key] = *jwtEntry.PublicKey
	log.Printf("%s - %s", jwtEntry.Key, data.JWT)
	log.Printf("API KEY GENERATED: from %s - %v", ctx.AccessId, data)
}

func (s *secretProviderImpl) RequiredScopes() []string {
	return []string{
		"api:example:user",
		"api:example:admin",
		"api:key:generator",
	}
}
