package mod

import "go/token"

type Package struct {
	PkgName    string
	PkgPath    string
	Targets    []*Func // funcs for fuzzing
	Funcs      []*Func
	Structs    []*Struct
	Interfaces []*Interface

	Fset *token.FileSet
}
