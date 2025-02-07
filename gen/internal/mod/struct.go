package mod

import "go/types"

const maxDepth int = 5

type Struct struct {
	StructName   string
	PkgName      string        // package name (should be the same as the package's package statement)
	PkgPath      string        // import path
	PkgDir       string        // local on-disk directory
	TypesStruct  *types.Struct // auxiliary information about a Struct from the go/types package
	TypesNamed   *types.Named  // auxiliary information about a Named from the go/types package
	Constructors []*Constructor

	AsPointer bool // detect that struct implement interface as pointer to this struct
}

func (s *Struct) String() string {
	return s.TypesNamed.String()
}

func (s *Struct) TypeString(qualifier types.Qualifier) string {
	return types.TypeString(s.TypesNamed, qualifier)
}

func getNotNativeTypes(s *types.Struct, depth int) map[types.Type]bool {
	notNatives := make(map[types.Type]bool)

	if depth == maxDepth {
		return notNatives
	}

	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)

		complex := getNamedOrSignatureTypes(field.Type())
		for _, notNative := range complex {
			switch named := notNative.(type) {
			case *types.Named:
				switch u := named.Underlying().(type) {
				case *types.Interface:
					notNatives[named] = true
				case *types.Struct:
					for nn, _ := range getNotNativeTypes(u, depth+1) {
						notNatives[nn] = true
					}
				case *types.Signature:
					notNatives[named] = true
				}
			case *types.Signature:
				notNatives[notNative] = true
			}
		}
	}

	return notNatives
}
