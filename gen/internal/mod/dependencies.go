package mod

import "go/types"

func (s *Struct) GetNotNativeTypes() map[*types.Var]bool {
	return getNotNativeTypes(s.TypesStruct, 1)
}

func HasNotNative(s *Struct) bool {
	return hasNotNative(s.TypesStruct, 1)
}

func GetInputParams(function *Func) []*types.Var {
	f := function.TypesFunc
	wrappedSig := f.Type().(*types.Signature)
	var inputParams []*types.Var

	for i := 0; i < wrappedSig.Params().Len(); i++ {
		v := wrappedSig.Params().At(i)
		inputParams = append(inputParams, v)
	}

	return inputParams
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
