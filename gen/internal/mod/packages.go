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

func (f *Func) GetSignature() *types.Signature {
	return f.TypesFunc.Type().(*types.Signature)
}

func (f *Func) TypeString(qualifier types.Qualifier) string {
	prefix := qualifier(f.TypesFunc.Pkg())
	if prefix == "" || f.GetSignature().Recv() != nil {
		return f.FuncName
	}

	return fmt.Sprintf("%s.%s", prefix, f.FuncName)
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

func (s Struct) String() string {
	return s.TypesNamed.String()
}

func (s Struct) TypeString(qualifier types.Qualifier) string {
	return types.TypeString(s.TypesNamed, qualifier)
}

type Interface struct {
	InterfaceName   string
	PkgName         string           // package name (should be the same as the package's package statement)
	PkgPath         string           // import path
	PkgDir          string           // local on-disk directory
	TypesInterface  *types.Interface // auxiliary information about a Interface from the go/types package
	TypesNamed      *types.Named     // auxiliary information about a Named from the go/types package
	Implementations []*Struct        // structs that implements this interface from explore packages
}

func (i Interface) String() string {
	return i.TypesNamed.String()
}

func (i Interface) TypeString(qualifier types.Qualifier) string {
	return types.TypeString(i.TypesNamed, qualifier)
}

type Package struct {
	PkgName      string
	PkgPath      string
	Targets      []*Func // funcs for fuzzing
	Constructors []*Func
	Funcs        []*Func
	Structs      []*Struct
	Interfaces   []*Interface
}

type PackagePattern struct {
	Funcs      []*Func
	Structs    []*Struct
	Interfaces []*Interface
}

func Merge(pkgs []*Package) *PackagePattern {
	pattern := PackagePattern{}

	for _, pkg := range pkgs {
		pattern.Funcs = append(pattern.Funcs, pkg.Funcs...)
		pattern.Structs = append(pattern.Structs, pkg.Structs...)
		pattern.Interfaces = append(pattern.Interfaces, pkg.Interfaces...)
	}

	return &pattern
}

func (pkg PackagePattern) GetSupportedInterfaces() map[string]*Interface {
	supportedInterfaces := make(map[string]*Interface)
	for _, pkgInterface := range pkg.Interfaces {
		if len(pkgInterface.Implementations) > 0 {
			supportedInterfaces[pkgInterface.String()] = pkgInterface
		}
	}

	return supportedInterfaces
}

func (pkg PackagePattern) GetSupportedStructs() map[string]*Struct {
	supportedStructs := make(map[string]*Struct)
	for _, pkgStruct := range pkg.Structs {
		supportedStructs[pkgStruct.String()] = pkgStruct
	}

	return supportedStructs
}
