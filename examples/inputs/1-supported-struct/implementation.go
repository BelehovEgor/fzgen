package supportedstruct

type A interface {
}

type F3 func(o int) int

type Complex struct {
	A  A
	F1 func(int) int
	F2 func(o int) int
	F3 F3
}

func F1(int) int {
	return 1
}

func Do(com Complex) int {
	return 1
}
