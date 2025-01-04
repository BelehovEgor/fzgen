package subpackage

import "fmt"

func MySubEmitter(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
