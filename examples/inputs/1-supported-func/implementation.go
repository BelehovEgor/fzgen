package supportedfunc

import (
	subpackage "github.com/thepudds/fzgen/examples/inputs/1-supported-func/sub-package"
)

type EmitFunc func(format string, args ...interface{})

func Log(emit EmitFunc, log string) {
	emit(log)
}

// func MyEmitter(format string, args ...interface{}) {
// 	fmt.Printf(format, args...)
// }

func Do() {
	Log(subpackage.MySubEmitter, "")
}
