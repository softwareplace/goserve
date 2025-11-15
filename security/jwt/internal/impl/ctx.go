package impl

import (
	goservectx "github.com/softwareplace/goserve/context"
)

func getDefaultCtx() *goservectx.DefaultContext {
	context := goservectx.NewDefaultCtx()
	context.SetRequesterId("gyo0V18QDj9Q1UWmZ2g7fc9sXrmlSthy3b8k9VO3MMv8dlEGtMtfIiPtJIUli0j")
	context.SetRoles("api:key:goserve-generator", "write:pets", "read:pets")
	return context
}
