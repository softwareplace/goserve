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
	"37c75552614a4eb58a2eb2d04928cdfd": "D/b7o5KGWe0SOF06r7bvKWyud95XVQwD9xp9NIDqMUWqt1xHz6PpIAF2jRo6pFGaaTwglXwql7QChU1fmQf7omQnjZImS9iWhKh9xvQEpXhygA5WAzBEPiekmyfH6LwkWgFQeFxi4spwX5J+m1LPMIrHZyjVqFOr01f3RaHAlBwxOwWdbQ0au32gVshGFY7Rt7d5RmMQATA0rQf0NGZlcIEM5ez8hBxjUHnKakGjYOITQsd570wvlFnRhvkvoxRfpAGAexXRAS8tImdiw/L7BVSbTKjwqSfweH59CK3JhHC/qdwDlSDA6rJWat4MOeb2qWbgbmlQV71QEFOZ9k78gdNz3FuFsMIQ4Swyf3dvBraTFlCjxDil7fIyTT1PJ8f8AvMcVdzWsXwWRl5+SgJvHcZI9nGmswzacRv2T008qUKm28m6By5Sd1ux38QghobBtpL2n3+lgEnov59/cStPHS4kSNrudeX1RtU7DPlqWZUyXkn4H+3tdlUXMufZcYekIkq3fIVsGHxRRGTRA1ILell9FBXwEVw/je2FsrzIZbPxZKnRb8WRbqNFreDf/9hdWLjKw4IaIddRUbGUSTLV3u94QbhDwsdFRmorMgKZd3yukVc=",
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
