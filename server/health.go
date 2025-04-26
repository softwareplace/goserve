package server

import (
	goservectx "github.com/softwareplace/goserve/context"
)

func (a *baseServer[T]) HealthResourceEnabled(value bool) Api[T] {
	a.healthResourceEnable = value
	return a
}

func (a *baseServer[T]) HealthResource() Api[T] {
	if a.healthResourceEnable {
		a.PublicRouter(a.healthHandler, "health", "GET")
	}
	return a
}

func (a *baseServer[T]) healthHandler(ctx *goservectx.Request[T]) {
	ctx.Ok(map[string]string{
		"status": "ok",
	})
}
