package supportedinterface

type Printer interface {
	Print(s string) (int, error)
	PrintPrinter() int
}

type MyPrinter struct {
}

func (printer MyPrinter) Print(s string) (int, error) {
	return len(s), nil
}

func (printer MyPrinter) PrintPrinter() int {
	return 1
}

func Do2() {
	Do(MyPrinter{})
}

func Do(printer Printer) {

}

type Emiter func(string) bool

func MyEmiter(s string) bool {
	return len(s) > 0
}

type StructWithComplexField struct {
	Printer Printer
	Emit    Emiter
}

func Target(s StructWithComplexField) {

}

func NewMyPrinter() MyPrinter {
	return MyPrinter{}
}
