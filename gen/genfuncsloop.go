package gen

import (
	"bytes"
	"errors"
	"fmt"
	"go/types"
	"io"
	"sort"

	"github.com/BelehovEgor/fzgen/gen/internal/mod"
)

var (
	errNoConstructorsMatch = errors.New("no matching constructor")
	errNoMethodsMatch      = errors.New("no methods found")
	errNoSteps             = errors.New("no supported methods found")
	errZeroFunctions       = errors.New("emitChainWrapper: zero functions")
	errSilentSkip          = errors.New("silently skipping wrapper generation")
)

// emitChainWrappers emits a set of fuzzing wrappers where possible for the list of functions passed in.
// Each wrapper consists of a target from a constructor and a set of steps that include invoking methods on the target.
// It might skip a function if it has no input parameters, or if it has a non-fuzzable parameter type.
func emitChainWrappers(
	outPkgPath string,
	pkgFuncs *mod.Package,
	typeContext *mod.TypeContext,
	wrapperPkgName string,
	options wrapperOptions,
) ([]byte, error) {
	possibleConstructors := typeContext.ValidConstructors
	if len(possibleConstructors) == 0 {
		return nil, errNoConstructorsMatch
	}

	// Build a map from the receiver type to a set of possible constructors
	// and possible steps with the same receiver type.
	type chain struct {
		recvType     string
		recv         *types.Named
		constructors []*mod.Constructor
		steps        []*mod.Func
	}

	recvTypes := make(map[string]*chain)
	for _, function := range pkgFuncs.Targets {
		// recvN will be the named type if the receiver is a pointer receiver.
		recvN := function.Receiver()
		if recvN == nil {
			continue
		}

		recvType := types.TypeString(recvN, nil)
		c := recvTypes[recvType]
		if c == nil {
			c = &chain{recvType: recvType, recv: recvN}
			recvTypes[recvType] = c
		}
		c.steps = append(c.steps, function)
	}

	if len(recvTypes) == 0 {
		return nil, errNoMethodsMatch
	}

	for constructor := range possibleConstructors {
		ctorResultN := constructor.ReturnType
		ctorType := types.TypeString(ctorResultN, nil)
		c := recvTypes[ctorType]
		if c == nil {
			// No methods found in loop above for this named type, so nothing to do with this possible constructor.
			continue
		}
		c.constructors = append(c.constructors, constructor)
	}

	// Put our chains in a deterministic order.
	var chains []*chain
	for _, v := range recvTypes {
		chains = append(chains, v)
	}
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].recvType < chains[j].recvType
	})

	// Prepare the output
	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	qualifier := mod.NewQualifier(pkgFuncs.PkgName, pkgFuncs.PkgPath, wrapperPkgName, outPkgPath, !options.qualifyAll)
	fabrics := mod.GenerateFabrics(pkgFuncs.Targets, typeContext, qualifier)
	init := mod.GenerateInitTestFunc(fabrics, typeContext, qualifier)

	// Emit the intro material
	emit("package %s\n\n", wrapperPkgName)
	emit(options.topComment)
	emit("import (\n")

	for _, importStr := range qualifier.GetImportStrings() {
		emit("\t%s\n", importStr)
	}

	emit(")\n\n")

	// Loop over our chains and emit fuzzing wrappers for each one.
	// We only return an error if all fail.
	var firstErr error
	var success bool
	for _, c := range chains {
		err := emitChainWrapper(emit, c.recv, c.steps, typeContext, qualifier, options)
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

// emitChainWrapper emits one fuzzing wrapper where possible for the list of functions passed in.
// It might skip a function if it has no input parameters, or if it has a non-fuzzable parameter
// type such as interface{}.
func emitChainWrapper(
	emit emitFunc,
	recv *types.Named,
	functions []*mod.Func,
	typeContext *mod.TypeContext,
	importQualifier *mod.ImportQualifier,
	options wrapperOptions,
) error {
	if len(functions) == 0 {
		return errZeroFunctions
	}

	variablesContext := mod.NewVariablesContext(importQualifier)

	// Determine our wrapper name, which includes the receiver's type if we are wrapping a method.
	wrapperName := fmt.Sprintf("Fuzz_Chain_%s", recv.Obj().Name())

	if !typeContext.IsSupported(recv.Obj().Type()) {
		recvTypeString := recv.Obj().Type().String()
		emit("// skipping %s because parameters include func, chan, or unsupported interface: %v\n\n", wrapperName, recvTypeString)
		return fmt.Errorf("%w: %s is unsupported", errUnsupportedParams, recvTypeString)
	}

	// Start emitting the wrapper function!
	// Start with the func declaration and the start of f.Fuzz.
	emit("func %s(f *testing.F) {\n", wrapperName)
	emit("\tf.Fuzz(func(t *testing.T, data []byte) {\n")

	recvName := variablesContext.CreateUniqueName("target")
	recvType := types.TypeString(recv, importQualifier.Qualifier)
	emit("\t\tvar %s %s\n", recvName, recvType)
	emit("\t\tfz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)\n")

	fillErrorName := variablesContext.CreateUniqueName("fillError")
	emit("\t\t%s := fz.Fill2(&%s)\n", fillErrorName, recvName)
	emit("\t\tif %s != nil", fillErrorName)
	_, ok := recv.Obj().Type().(*types.Pointer)
	if ok {
		emit(" || %s == nil", recvName)
	}
	emit(" {\n")
	emit("\t\t\treturn\n")
	emit("\t\t}\n")

	// put our functions we want to wrap into a deterministic order
	sort.Slice(functions, func(i, j int) bool {
		// types.Func.String outputs strings like:
		//   func (github.com/thepudds/fzgo/genfuzzfuncs/examples/test-constructor-injection.A).ValMethodWithArg(i int) bool
		// works ok for clustering results, though pointer receiver and non-pointer receiver methods don't cluster.
		// could strip '*' or sort another way, but probably ok, at least for now.
		return functions[i].TypesFunc.String() < functions[j].TypesFunc.String()
	})

	emit("\tsteps := []fuzzer.Step{\n")

	// loop over our the functions we are wrapping, emitting a wrapper where possible.
	var emittedSteps int
	for _, function := range functions {
		err := emitChainStep(emit, wrapperName, recvName, function, typeContext, importQualifier)
		if errors.Is(err, errSilentSkip) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error processing %s: %v", function.FuncName, err)
		}
		emittedSteps++
	}
	// close out steps slice
	emit("\t}\n\n")

	if emittedSteps == 0 {
		// TODO: we could handle this better, but let's close out this wrapper in case there is another
		// chain that is useful. The whole output file will be skipped if this was the only candidate chain.
		emit("\t\t_, _, _ = fz, target, steps")
		// close out the f.Fuzz func
		emit("\t})\n")
		// close out test func
		emit("}\n\n")
		return errNoSteps
	}

	// emit the chain func
	emit("\t// Execute a specific chain of steps, with the count, sequence and arguments controlled by fz.Chain\n")
	if !options.parallel {
		emit("\tfz.Chain(steps)\n")
	} else {
		emit("\tfz.Chain(steps, fuzzer.ChainParallel)\n")
	}

	// close out the f.Fuzz func
	emit("\t})\n")

	// close out test func
	emit("}\n\n")

	return nil
}

// emitChainStep emits one fuzzing step if possible.
// It takes a list of possible constructors to insert into the step body if the
// constructor is suitable for creating the receiver of a wrapped method.
// qualifyAll indicates if all variables should be qualified with their package.
func emitChainStep(
	emit emitFunc,
	wrapperName, recvName string,
	function *mod.Func,
	typeContext *mod.TypeContext,
	importQualifier *mod.ImportQualifier,
) error {
	// new context in lambda
	variableContext := mod.NewVariablesContext(importQualifier)

	wrappedSig := function.GetSignature()

	// Start building up our list of parameters we will use in input
	// parameters to the new wrapper func we are about to emit.
	var inputParams []*types.Var

	// Also add in the parameters for the function under test.
	// Note that we allow zero len Params because we only get this far
	// if we have a match on the receiver with the target object.
	for i := 0; i < wrappedSig.Params().Len(); i++ {
		v := wrappedSig.Params().At(i)
		inputParams = append(inputParams, v)
	}

	paramReprs := make([]*paramRepr, 0, len(inputParams)+1)
	for _, v := range inputParams {
		paramReprs = append(paramReprs, newParam(v, variableContext, importQualifier.Qualifier))
	}

	// Check if we have an interface or function pointer in our desired parameters,
	// which we can't fill with values during fuzzing.
	isSupported := true
	var unsupportedParam *types.Var
	for _, v := range inputParams {
		if !typeContext.IsSupported(v.Type()) {
			isSupported = false
			unsupportedParam = v
			break
		}
	}
	if !isSupported {
		unsupportedParamType := unsupportedParam.Type().String()
		emit("// skipping %s because parameters include unsupported type: %v\n\n", wrapperName, unsupportedParamType)
		return errSilentSkip
	}

	// Start emitting the wrapper function, inside of a fzgen/fuzzer.Step. Will be similar to:
	//   Step{
	// 	   Name: "input int",
	// 	   Func: func(a int) int {
	//	          return a
	//     }
	//   },
	emit("\t{\n")
	emit("\t\tName: \"%s\",\n", wrapperName)
	emit("\t\tFunc: func(")

	// For independent wrappers, in some cases we need to emit fz.Fill to handle
	// creating rich args that are beyond what cmd/go can fuzz (including because
	// the standard go test infrastructure will be calling the wrappers we created).
	// In contrast, for the chain steps we are creating here, we emit the same
	// code for both nativeSupport and fillRequired, and handle the difference at run time.
	// This is because for chain steps, we never emit fz.Fill calls because at run time
	// we are the ones to call the function literal, and hence we create those rich args
	// at run time, and hence here we just create function literals with the arguments
	// we want.
	//
	// The result for this line will end up similar to:
	//    Func: func(s string, i *int) {
	// Iterate over the our input parameters and emit.
	for i, p := range paramReprs {
		// want: foo string, bar int
		if i > 0 {
			// need a comma if something has already been emitted
			emit(", ")
		}
		emit("%s %s", p.paramName, p.typ)
	}
	emit(") ")

	// TODO: centralize error check logic
	results := wrappedSig.Results()
	if results.Len() > 0 && !(results.Len() == 1 && results.At(0).Type().String() == "error") {
		emit("(") // goimports should clean up paren if it is not needed
		for i := 0; i < results.Len(); i++ {
			if i > 0 {
				emit(", ")
			}
			returnTypeStringWithSelector := types.TypeString(results.At(i).Type(), importQualifier.Qualifier)
			emit(returnTypeStringWithSelector)
		}
		emit(")")
	}
	emit(" {\n")

	// Emit the call to the wrapped function.
	// collisionOffset is 0 because we do not have a constructor within this function
	// literal we are creating and hence we don't need to worry about calculating
	// a collisionOffset.
	// Include a 'return' if we have a non-error return value for our wrapped func.
	if results.Len() > 0 && !(results.Len() == 1 && results.At(0).Type().String() == "error") {
		emit("\treturn ")
	}
	emitWrappedFunc(emit, function, &paramRepr{paramName: recvName}, paramReprs, importQualifier.Qualifier)

	// close out the func as well as the Step struct
	emit("\t\t},\n")
	emit("\t},\n")
	return nil
}
