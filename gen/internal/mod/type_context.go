package mod

import (
	"fmt"
	"go/types"
	"regexp"
)

var supportedInterfaces = map[string]bool{
	"io.Writer":       true,
	"io.Reader":       true,
	"io.ReaderAt":     true,
	"io.WriterTo":     true,
	"io.Seeker":       true,
	"io.ByteScanner":  true,
	"io.RuneScanner":  true,
	"io.ReadSeeker":   true,
	"io.ByteReader":   true,
	"io.RuneReader":   true,
	"io.ByteWriter":   true,
	"io.ReadWriter":   true,
	"io.ReaderFrom":   true,
	"io.StringWriter": true,
	"io.Closer":       true,
	"io.ReadCloser":   true,
	"context.Context": true,
	"error":           true,
	"interface{}":     true,
}

type TypeContext struct {
	ExistedFuncs        []*Func
	ExistedStructs      []*Struct
	ExistedInterfaces   []*Interface
	SupportedInterfaces map[*Interface]bool

	ValidConstructors map[*Constructor]bool
	canDo             map[string]bool

	required      map[string][]*Constructor
	requiredCount map[*Constructor]int

	constructorPattern *regexp.Regexp
}

func NewTypeContext(constructorPattern *regexp.Regexp) *TypeContext {
	canDo := make(map[string]bool)
	for key, value := range supportedInterfaces {
		canDo[key] = value
	}

	return &TypeContext{
		ValidConstructors:   make(map[*Constructor]bool),
		SupportedInterfaces: make(map[*Interface]bool),
		canDo:               canDo,
		required:            make(map[string][]*Constructor),
		requiredCount:       make(map[*Constructor]int),
		constructorPattern:  constructorPattern,
	}
}

func (tc *TypeContext) AddType(t types.Type) {
	typeName := t.String()
	tc.canDo[typeName] = true
}

func (tc *TypeContext) AddFunc(f *Func) {
	tc.ExistedFuncs = append(tc.ExistedFuncs, f)
}

func (tc *TypeContext) AddStruct(s *Struct) {
	tc.ExistedStructs = append(tc.ExistedStructs, s)

	for _, i := range tc.ExistedInterfaces {
		asInterface := castStructAsInterfaceImplementation(s, i)
		if asInterface != nil {
			tc.setInterfaceImplementation(i, asInterface)
		}
	}
}

func (tc *TypeContext) AddInterface(i *Interface) {
	tc.ExistedInterfaces = append(tc.ExistedInterfaces, i)

	for _, s := range tc.ExistedStructs {
		asInterface := castStructAsInterfaceImplementation(s, i)
		if asInterface != nil {
			tc.setInterfaceImplementation(i, asInterface)
		}
	}
}

func (tc *TypeContext) TryAddAsConstructor(f *Func) {
	con := tryCastToConstructor(f)

	if con == nil || !tc.constructorPattern.MatchString(con.Func.FuncName) {
		return
	}

	params := con.Func.Params()

	var required []types.Type
	hasUnsupported := false
	for _, param := range params {
		if !tc.IsSupported(param.Type()) {
			required = append(required, param.Type())
			hasUnsupported = true
		}
	}

	if !hasUnsupported {
		tc.addValidConstructor(con)

		for _, constructorWithUnsupportedParam := range tc.required[con.ReturnType.String()] {
			tc.requiredCount[constructorWithUnsupportedParam]--

			if tc.requiredCount[constructorWithUnsupportedParam] == 0 {
				delete(tc.requiredCount, constructorWithUnsupportedParam)
				tc.addValidConstructor(constructorWithUnsupportedParam)
			}
		}

		delete(tc.required, con.ReturnType.String())
	} else {
		tc.requiredCount[con] = len(required)
		for _, param := range required {
			tc.required[param.String()] = append(tc.required[param.String()], con)
		}
	}
}

func (tc *TypeContext) IsSupported(
	t types.Type,
) bool {
	// Switch to check if we might be able to fill this type.
	switch u := t.Underlying().(type) {
	case *types.Pointer:
		if !tc.IsSupported(u.Elem()) {
			return false
		}
	case *types.Slice:
		if !tc.IsSupported(u.Elem()) {
			return false
		}
	case *types.Array:
		if !tc.IsSupported(u.Elem()) {
			return false
		}
	case *types.Map:
		if !tc.IsSupported(u.Elem()) {
			return false
		}
	case *types.Struct:
		return true
	case *types.Interface:
		typeName := t.String()
		return tc.canDo[typeName]
	case *types.Signature:
		typeName := t.String()
		return tc.canDo[typeName]
	case *types.Chan:
		return false
	case *types.Basic:
		switch u.Kind() {
		case types.Uintptr, types.UnsafePointer, types.Complex64, types.Complex128:
			return true
		}
	}

	return true
}

// is implements as value
// is implements as ptr
func castStructAsInterfaceImplementation(_struct *Struct, _interface *Interface) *Struct {
	if types.AssignableTo(_struct.TypesNamed, _interface.TypesInterface) {
		return _struct
	} else if types.AssignableTo(types.NewPointer(_struct.TypesNamed), _interface.TypesInterface) {
		return &Struct{
			StructName:  _struct.StructName,
			PkgName:     _struct.PkgName,
			PkgPath:     _struct.PkgPath,
			PkgDir:      _struct.PkgDir,
			TypesStruct: _struct.TypesStruct,
			TypesNamed:  _struct.TypesNamed,
			AsPointer:   true,
		}
	}

	return nil
}

func isConstructorOf(named *types.Named, possibleCtor *Constructor) bool {
	return types.TypeString(named, nil) == types.TypeString(possibleCtor.ReturnType, nil)
}

func (tc *TypeContext) setInterfaceConstructors(constructor *Constructor) {
	for _, i := range tc.ExistedInterfaces {
		if isConstructorOf(i.TypesNamed, constructor) {
			tc.setInterfaceConstructor(i, constructor)
		}
	}
}

func (tc *TypeContext) setStructConstructors(constructor *Constructor) {
	for _, s := range tc.ExistedStructs {
		if isConstructorOf(s.TypesNamed, constructor) {
			s.Constructors = append(s.Constructors, constructor)
		}
	}
}

func (tc *TypeContext) setInterfaceConstructor(i *Interface, constructor *Constructor) {
	tc.SupportedInterfaces[i] = true
	i.Constructors = append(i.Constructors, constructor)
}

func (tc *TypeContext) setInterfaceImplementation(i *Interface, s *Struct) {
	tc.SupportedInterfaces[i] = true
	i.Implementations = append(i.Implementations, s)
}

func (tc *TypeContext) addValidConstructor(constructor *Constructor) {
	tc.ValidConstructors[constructor] = true
	constructor.IsSupported = true

	tc.setInterfaceConstructors(constructor)
	tc.setStructConstructors(constructor)

	typeName := constructor.ReturnType.String()
	tc.canDo[typeName] = true
}

func tryCastToConstructor(f *Func) *Constructor {
	if f == nil || !f.TypesFunc.Exported() {
		return nil
	}

	ctorSig, ok := f.TypesFunc.Type().(*types.Signature)
	if !ok || ctorSig.Recv() != nil {
		return nil
	}

	ctorResults := ctorSig.Results()
	if ctorResults.Len() > 2 || ctorResults.Len() == 0 {
		return nil
	}

	secondResultIsErr := false
	if ctorResults.Len() == 2 {
		// We allow error type as second return value
		secondResult := ctorResults.At(1)
		_, ok = secondResult.Type().Underlying().(*types.Interface)
		if ok && secondResult.Type().String() == "error" {
			secondResultIsErr = true
		} else {
			return nil
		}
	}

	ctorResult := ctorResults.At(0)
	ctorResultN, err := namedType(ctorResult)
	if err != nil {
		return nil
	}

	return &Constructor{
		Func:              f,
		ReturnType:        ctorResultN,
		SecondResultIsErr: secondResultIsErr,
	}
}

func namedType(recv *types.Var) (*types.Named, error) {
	reportErr := func() (*types.Named, error) {
		return nil, fmt.Errorf("expected pointer or named type: %+v", recv.Type())
	}

	switch t := recv.Type().(type) {
	case *types.Pointer:
		if t.Elem() == nil {
			return reportErr()
		}
		n, ok := t.Elem().(*types.Named)
		if ok {
			return n, nil
		}
	case *types.Named:
		return t, nil
	}
	return reportErr()
}
