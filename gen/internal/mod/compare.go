package mod

import "go/types"

func FindSuitables(s *types.Signature, existedFuncs []*Func) []*Func {
	var sameFunc []*Func
	for _, f := range existedFuncs {
		if compareSignatures(s, f.GetSignature()) {
			sameFunc = append(sameFunc, f)
		}
	}

	if sameFunc == nil {
		return nil
	}

	return sameFunc
}

func compareSignatures(sig1, sig2 *types.Signature) bool {
	if sig1.Params().Len() != sig2.Params().Len() {
		return false
	}

	for i := 0; i < sig1.Params().Len(); i++ {
		param1 := sig1.Params().At(i)
		param2 := sig2.Params().At(i)

		if !types.Identical(param1.Type(), param2.Type()) {
			return false
		}
	}

	if sig1.Results().Len() != sig2.Results().Len() {
		return false
	}

	for i := 0; i < sig1.Results().Len(); i++ {
		result1 := sig1.Results().At(i)
		result2 := sig2.Results().At(i)

		if !types.Identical(result1.Type(), result2.Type()) {
			return false
		}
	}

	return true
}
