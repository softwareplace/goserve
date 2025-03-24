package error

import (
	"fmt"
	apicontext "github.com/softwareplace/goserve/context"
	"runtime"
)

func Handler(try func(), catch func(err error)) {
	defer func() {
		if r := recover(); r != nil {
			_, file, line, ok := runtime.Caller(2) // Adjust caller depth to log where the error originates
			var errMessage = fmt.Sprintf("panic occurred at %s:%d - %v", file, line, r)

			if !ok {
				errMessage = fmt.Sprintf("panic occurred - %v", r)
			}

			catch(Wrapper(fmt.Errorf(errMessage), "Recovered panic"))
		}
	}()
	try()
}

func Wrapper(err error, message string) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("%s: %w", message, err)
	}
	return fmt.Errorf("%s (%s:%d): %w", message, file, line, err)
}

type ApiHandler[T apicontext.Principal] interface {
	Handler(ctx *apicontext.Request[T], err error, source string)
}
