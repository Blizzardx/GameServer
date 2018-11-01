package Common

import (
	"fmt"
	"runtime/debug"
)

func SafeCall(f func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf(string(debug.Stack()))
		}
	}()
	f()
}
func SafeCallWithCrashCallback(f func(), onCrash func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf(string(debug.Stack()))
			onCrash()
		}
	}()
	f()
}
