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

type Constructor struct {
	Func *Func

	ReturnType        *types.Named
	SecondResultIsErr bool

	IsSupported bool
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

func (f *Func) Params() []*types.Var {
	wrappedSig := f.GetSignature()
	var inputParams []*types.Var
	for i := 0; i < wrappedSig.Params().Len(); i++ {
		v := wrappedSig.Params().At(i)
		inputParams = append(inputParams, v)
	}
	return inputParams
}
