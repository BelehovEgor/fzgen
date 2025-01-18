package mod

import "go/types"

type Interface struct {
	InterfaceName   string
	PkgName         string           // package name (should be the same as the package's package statement)
	PkgPath         string           // import path
	PkgDir          string           // local on-disk directory
	TypesInterface  *types.Interface // auxiliary information about a Interface from the go/types package
	TypesNamed      *types.Named     // auxiliary information about a Named from the go/types package
	Implementations []*Struct        // structs that implements this interface from explore packages
	Constructors    []*Constructor
}

func (i *Interface) String() string {
	return i.TypesNamed.String()
}

func (i *Interface) TypeString(qualifier types.Qualifier) string {
	return types.TypeString(i.TypesNamed, qualifier)
}

func (i *Interface) IsSupported() bool {
	return len(i.Implementations) > 0 || len(i.Constructors) > 0
}
