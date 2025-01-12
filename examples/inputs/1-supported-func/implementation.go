package supportedfunc

import (
	"fmt"

	subpackage "github.com/thepudds/fzgen/examples/inputs/1-supported-func/sub-package"
)

type EmitFunc func(lvl int, format string, args ...interface{})

func Log(emit EmitFunc, lvl int, log string) {
	emit(lvl, log)
}

func MyEmitter(lvl int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func Do() {
	Log(subpackage.MySubEmitter, 12, "")
}