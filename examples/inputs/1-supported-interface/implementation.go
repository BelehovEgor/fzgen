package supportedinterface

import "fmt"

type Printer interface {
	Print(s string) int
}

type MyPrinter struct {
}

func (printer MyPrinter) Print(s string) int {
	return len(s)
}

type B struct {
	I, J int
	S    string
}

type MyPointerPrinter struct {
	B        B
	I        int
	Printer  Printer
	Printer_ Printer
	Printer2 MyPrinter
}

func (printer *MyPointerPrinter) Print(s string) int {
	return len(s)
}

type IInterface interface {
	Do(printer Printer) bool
}

type A struct {
}

func (a A) Do(printer Printer) bool {
	return true
}

func MySubEmitter(lvl int, format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
