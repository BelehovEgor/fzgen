package mod

import (
	"fmt"
	"go/types"
)

// Func represents a discovered function that will be fuzzed.
type Func struct {
	FuncName  string
	PkgName   string      // package name (should be the same as the package's package statement)
	PkgPath   string      // import path
	PkgDir    string      // local on-disk directory
	TypesFunc *types.Func // auxiliary information about a Func from the go/types package
}

// FuzzName returns the '<pkg>.<OrigFuzzFunc>' string.
// For example, it might be 'fmt.FuzzFmt'. In fzgo,
// this was used in messages, and as part of the path when creating
// the corpus location under testdata.
func (f *Func) FuzzName() string {
	return fmt.Sprintf("%s.%s", f.PkgName, f.FuncName)
}

func (f *Func) String() string {
	return f.FuzzName()
}

func (i Func) GetFullName() string {
	return i.PkgPath + i.FuncName
}

type Struct struct {
	StructName  string
	PkgName     string        // package name (should be the same as the package's package statement)
	PkgPath     string        // import path
	PkgDir      string        // local on-disk directory
	TypesStruct *types.Struct // auxiliary information about a Struct from the go/types package
	TypesNamed  *types.Named  // auxiliary information about a Named from the go/types package

	AsPointer bool // detect that struct implement interface as pointer to this struct
}

type Interface struct {
	InterfaceName   string
	PkgName         string           // package name (should be the same as the package's package statement)
	PkgPath         string           // import path
	PkgDir          string           // local on-disk directory
	TypesInterface  *types.Interface // auxiliary information about a Struct from the go/types package
	Implementations []Struct         // structs that implements this interface from explore packages
}

func (i Interface) GetFullName() string {
	return i.PkgPath + i.InterfaceName
}

// TODO: probably rename internal/mod to internal/pkg.
type Package struct {
	PkgName      string
	PkgPath      string
	Targets      []Func // funcs for fuzzing
	Constructors []Func
	Funcs        []Func
	Structs      []Struct
	Interfaces   []Interface
}

func (pkg Package) GetSupportedInterfaces() map[string]Interface {
	supportedInterfaces := make(map[string]Interface)
	for _, pkgInterface := range pkg.Interfaces {
		if len(pkgInterface.Implementations) > 0 {
			supportedInterfaces[pkgInterface.GetFullName()] = pkgInterface
		}
	}

	return supportedInterfaces
}

func (pkg Package) GetExistedFunctions() map[string]Func {
	existsFunctions := make(map[string]Func)
	for _, pkgFunc := range pkg.Funcs {
		sig := pkgFunc.TypesFunc.Type().Underlying().(*types.Signature)

		existsFunctions[sig.String()] = pkgFunc
	}

	return existsFunctions
}
