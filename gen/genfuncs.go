package gen

import (
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"io"
	"os"
	"sort"

	"github.com/thepudds/fzgen/gen/internal/mod"
)

type wrapperOptions struct {
	qualifyAll bool   // qualify all variables with package name
	topComment string // additional comment for top of generated file.
}

type emitFunc func(format string, args ...interface{})

var (
	errNoFunctionsMatch  = errors.New("no fuzzable functions found")
	errUnsupportedParams = errors.New("unsupported parameters")
)

// emitIndependentWrappers emits fuzzing wrappers where possible for the list of functions passed in.
// It might skip a function if it has no input parameters, or if it has a non-fuzzable parameter
// type such as interface{}.
func emitIndependentWrappers(
	outPkgPath string,
	pkgFuncs *mod.Package,
	typeContext *mod.TypeContext,
	wrapperPkgName string,
	options wrapperOptions,
) ([]byte, error) {
	if len(pkgFuncs.Targets) == 0 {
		return nil, fmt.Errorf("%w: 0 matching functions", errNoFunctionsMatch)
	}

	// prepare the output
	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	// put our functions we want to wrap into a deterministic order
	sort.Slice(pkgFuncs.Targets, func(i, j int) bool {
		// types.Func.String outputs strings like:
		//   func (github.com/thepudds/fzgo/genfuzzfuncs/examples/test-constructor-injection.A).ValMethodWithArg(i int) bool
		// works ok for clustering results, though pointer receiver and non-pointer receiver methods don't cluster.
		// could strip '*' or sort another way, but probably ok, at least for now.
		return pkgFuncs.Targets[i].TypesFunc.String() < pkgFuncs.Targets[j].TypesFunc.String()
	})

	qualifier := mod.NewQualifier(pkgFuncs.PkgName, pkgFuncs.PkgPath, wrapperPkgName, outPkgPath, !options.qualifyAll)
	fabrics := mod.GenerateFabrics(pkgFuncs.Targets, typeContext, qualifier)
	init := mod.GenerateInitTestFunc(fabrics, typeContext, qualifier)

	// emit the intro material
	emit("package %s\n\n", wrapperPkgName)
	emit(options.topComment)
	emit("import (\n")

	for _, importStr := range qualifier.GetImportStrings() {
		emit("\t%s\n", importStr)
	}

	emit(")\n\n")

	// Loop over our the functions we are wrapping, emitting a wrapper where possible.
	// We only return an error if all fail.
	var firstErr error
	var success bool
	for _, function := range pkgFuncs.Targets {
		err := emitIndependentWrapper(
			emit, *function, typeContext, qualifier)
		if err != nil && firstErr == nil {
			firstErr = err
		}
		if err == nil {
			success = true
		}
	}
	if !success {
		return nil, firstErr
	}

	funcs := make([]*mod.GeneratedFunc, 0, len(fabrics))
	for _, value := range fabrics {
		funcs = append(funcs, value)
	}
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Name > funcs[j].Name
	})

	for _, value := range funcs {
		emit("%s\n", value.Body)
	}

	emit("%s\n\n", init.Body)

	return buf.Bytes(), nil
}

// paramRepr contains string representations of inputParams to the wrapper function that we are
// creating. It includes params for the function under test, as well as in some cases
// args for a related constructor.
type paramRepr struct {
	paramName string
	typ       string
	v         *types.Var
}

func newParam(v *types.Var, varContext *mod.VariablesContext, qualifier types.Qualifier) *paramRepr {
	if v == nil {
		return nil
	}

	typeStringWithSelector := types.TypeString(v.Type(), qualifier)
	paramName := varContext.CreateUniqueNameForVariable(v)
	return &paramRepr{paramName: paramName, typ: typeStringWithSelector, v: v}
}

// emitIndependentWrapper emits one fuzzing wrapper if possible.
// It takes a list of possible constructors to insert into the wrapper body if the
// constructor is suitable for creating the receiver of a wrapped method.
// qualifyAll indicates if all variables should be qualified with their package.
func emitIndependentWrapper(
	emit emitFunc,
	function mod.Func,
	typeContext *mod.TypeContext,
	qualifier *mod.ImportQualifier) error {
	f := function.TypesFunc
	wrappedSig, ok := f.Type().(*types.Signature)
	if !ok {
		return fmt.Errorf("function %s is not *types.Signature (%+v)", function, f)
	}

	// Get our receiver, which might be nil if we don't have a receiver
	recv := wrappedSig.Recv()

	// Determine our wrapper name, which includes the receiver's type if we are wrapping a method.
	var wrapperName string
	if recv == nil {
		wrapperName = fmt.Sprintf("Fuzz_%s", f.Name())
	} else {
		n, err := namedType(recv)
		if err != nil {
			// output to stderr, but don't treat as fatal error.
			fmt.Fprintf(os.Stderr, "genfuzzfuncs: warning: createWrapper: failed to determine receiver type: %v: %v\n", recv, err)
			return nil
		}
		recvNamedTypeLocalName := n.Obj().Name()
		wrapperName = fmt.Sprintf("Fuzz_%s_%s", recvNamedTypeLocalName, f.Name())
	}

	varContext := mod.NewVariablesContext(qualifier)

	var inputParams []*types.Var
	for i := 0; i < wrappedSig.Params().Len(); i++ {
		v := wrappedSig.Params().At(i)
		inputParams = append(inputParams, v)
	}
	if len(inputParams) == 0 {
		// skip this wrapper, not useful for fuzzing if no inputs (no receiver, no parameters).
		return fmt.Errorf("%w: %s has 0 input params", errNoFunctionsMatch, function.FuncName)
	}

	paramReprs := make([]*paramRepr, 0, len(inputParams)+1)

	var recvParam *paramRepr
	if recv != nil {
		recvParam = newParam(recv, varContext, qualifier.Qualifier)
		paramReprs = append(paramReprs, recvParam)
	}

	for _, v := range inputParams {
		paramReprs = append(paramReprs, newParam(v, varContext, qualifier.Qualifier))
	}

	// Check if we have an interface or function pointer in our desired parameters,
	// which we can't fill with values during fuzzing.
	isSupported := true
	var unsupportedParam *types.Var
	for _, v := range inputParams {
		if v.Type().String() == "func(net/http.ResponseWriter, *net/http.Request)" {
			i := 1
			i++
		}

		if !typeContext.IsSupported(v.Type()) {
			isSupported = false
			unsupportedParam = v
			break
		}
	}
	if !isSupported {
		unsupportedParamType := unsupportedParam.Type().String()
		emit("// skipping %s because parameters include unsupported type: %v\n\n", wrapperName, unsupportedParamType)
		return fmt.Errorf("%w: %s takes %s", errUnsupportedParams, function.FuncName, unsupportedParamType)
	}

	// Start emitting func
	emit("func %s(f *testing.F) {\n", wrapperName)
	emit("\tf.Fuzz(func(t *testing.T, data []byte) {\n")
	for _, p := range paramReprs {
		emit("\t\tvar %s %s\n", p.paramName, p.typ)
	}

	fillErrVarName := varContext.CreateUniqueName("err")
	emit("\t\tfz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)\n")
	emit("\t\t%s := fz.Fill2(", fillErrVarName)
	for i, p := range paramReprs {
		if i > 0 {
			emit(", ")
		}
		emit("&%s", p.paramName)
	}
	emit(")\n")

	emitFillResultCheck(emit, fillErrVarName, paramReprs)

	emit("\n")

	args := paramReprs
	if recvParam != nil {
		args = paramReprs[1:]
	}

	emitWrappedFunc(emit, function, recvParam, args, qualifier.Qualifier)
	emit("\t})\n")
	emit("}\n\n")

	return nil
}

func emitFillResultCheck(emit emitFunc, fillErrorName string, allParams []*paramRepr) {
	emit("\tif %s != nil", fillErrorName)

	for _, p := range allParams {
		_, ok := p.v.Type().(*types.Pointer)
		if ok {
			emit(" || ")
			paramName := p.paramName
			emit("%s == nil", paramName)
		}
	}
	emit(" {\n")
	emit("\t\treturn\n")
	emit("\t}\n")
}

// emitWrappedFunc emits the call to the function under test.
// A target that is not "" indicates the caller wants to use a
// specific target name in place of any receiver name.
// For example, a target set to "target" would result in "target.Load(key)".
func emitWrappedFunc(
	emit emitFunc,
	f mod.Func,
	recv *paramRepr,
	paramReprs []*paramRepr,
	qualifier types.Qualifier,
) {
	switch {
	case recv != nil:
		emit("\t%s.%s(", recv.paramName, f.TypeString(qualifier))
	default:
		emit("\t%s(", f.TypeString(qualifier))
	}
	emitArgs(emit, f, paramReprs)
	emit(")\n")
}

// emitArgs emits the arguments needed to call a signature, including handling renaming arguments
// based on collisions with package name or other parameters.
func emitArgs(
	emit emitFunc,
	f mod.Func,
	paramReprs []*paramRepr,
) {
	sig := f.GetSignature()
	for i := 0; i < sig.Params().Len(); i++ {
		paramName := paramReprs[i].paramName
		if i > 0 {
			emit(", ")
		}
		emit(paramName)
	}
	if sig.Variadic() {
		// last argument needs an elipsis
		emit("...")
	}
}

// namedType returns a *types.Named if the passed in
// *types.Var is a *types.Pointer or a *types.Named.
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
