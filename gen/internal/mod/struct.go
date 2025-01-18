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

func (s *Struct) GetNotNativeTypes() map[*types.Var]bool {
	return getNotNativeTypes(s.TypesStruct, 1)
}

func getNotNativeTypes(s *types.Struct, depth int) map[*types.Var]bool {
	if depth > maxDepth {
		return nil
	}

	notNative := make(map[*types.Var]bool)
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		switch u := field.Type().Underlying().(type) {
		case *types.Interface, *types.Signature:
			notNative[field] = true
		case *types.Struct:
			fieldNotNative := getNotNativeTypes(u, depth+1)
			if fieldNotNative == nil {
				continue
			}

			for v := range fieldNotNative {
				notNative[v] = true
			}
		}
	}

	return notNative
}
