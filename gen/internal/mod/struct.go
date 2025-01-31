package mod

import "go/types"

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

func (s *Struct) GetNotNativeTypes() map[types.Type]bool {
	return getNotNativeTypes(s.TypesStruct, 1)
}

func getNotNativeTypes(s *types.Struct, depth int) map[types.Type]bool {
	if depth > maxDepth {
		return nil
	}

	notNatives := make(map[types.Type]bool)
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)

		complex := getNamedOrSignatureTypes(field.Type())
		if len(complex) > 0 {
			for _, notNative := range complex {
				notNatives[notNative] = true
			}
		}
	}

	return notNatives
}
