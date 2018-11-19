package ws

import (
	"fmt"
	"runtime"
)

// NotYetImplemented returns a nice formatted error for unimplemented API function call
func NotYetImplemented() error {
	fpc, _, _, ok := runtime.Caller(1)
	if ok == false {
		panic("Could not found caller")
	}
	fun := runtime.FuncForPC(fpc)

	return fmt.Errorf("%s is not yet implemented", fun.Name())
}
