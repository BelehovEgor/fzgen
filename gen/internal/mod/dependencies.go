package mod

import "go/types"

func (s *Struct) GetNotNativeTypes() []*types.Var {
	notNative := make([]*types.Var, 0)
	for i := 0; i < s.TypesStruct.NumFields(); i++ {
		field := s.TypesStruct.Field(i)
		switch u := field.Type().Underlying().(type) {
		case *types.Interface, *types.Signature:
			notNative = append(notNative, field)
		case *types.Struct:
			if hasNotNative(u, 1) {
				notNative = append(notNative, field)
			}
		}
	}

	return notNative
}

func HasNotNative(s *Struct) bool {
	return hasNotNative(s.TypesStruct, 1)
}

func hasNotNative(s *types.Struct, depth int) bool {
	if depth > maxDepth {
		return false
	}

	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		switch field.Type().Underlying().(type) {
		case *types.Interface, *types.Signature:
			return true
		case *types.Struct:
			has := hasNotNative(s, depth+1)
			if has {
				return true
			}
		}
	}

	return false
}
