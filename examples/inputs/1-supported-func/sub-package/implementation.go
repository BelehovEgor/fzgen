package subpackage

import "fmt"

func MySubEmitter(lvl int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
