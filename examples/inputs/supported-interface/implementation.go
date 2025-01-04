package supportedinterface

type Printer interface {
	Print(s string) int
}

type MyPrinter struct {
}

func (printer MyPrinter) Print(s string) int {
	return len(s)
}

type MyPointerPrinter struct {
}

func (printer *MyPointerPrinter) Print(s string) int {
	return len(s)
}

type IInterface interface {
	Do(x int) bool
}

type A struct {
}

func (a A) Do(x int) bool {
	return x%2 == 0
}
