package principal

import "github.com/softwareplace/http-utils/api_context"

type PService[T api_context.ApiContextData] interface {
	Roles(ctx api_context.ApiRequestContext[T]) []string
}
