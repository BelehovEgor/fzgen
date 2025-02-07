package mod

import (
	"fmt"
	"go/types"
)

func GenerateMockFabrics(
	targets []*Func,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) ([]*GeneratedFunc, []byte) {
	var typesThatNeededMock []types.Type

	for _, target := range targets {
		for _, param := range target.Params() {
			baseTypes := getNamedOrSignatureTypes(param.Type())
			for _, baseType := range baseTypes {
				switch t := baseType.(type) {
				case *types.Named:
					switch u := t.Underlying().(type) {
					case *types.Interface:
						typesThatNeededMock = append(typesThatNeededMock, t)
					case *types.Struct:
						for notNativeType := range getNotNativeTypes(u, 1) {
							switch structNotNative := notNativeType.(type) {
							case *types.Named:
								if types.IsInterface(structNotNative) {
									typesThatNeededMock = append(typesThatNeededMock, structNotNative)
								}
							case *types.Signature:
								typesThatNeededMock = append(typesThatNeededMock, structNotNative)
							}
						}
					case *types.Signature:
						typesThatNeededMock = append(typesThatNeededMock, t)
					}
				case *types.Signature:
					typesThatNeededMock = append(typesThatNeededMock, t)
				}
			}
		}
	}

	mockFabrics := createMocks(typesThatNeededMock, typeContext, qualifier, relativePackagePath, maxDepth)

	return mockFabrics, generateMockeryYaml(typesThatNeededMock)
}

func generateMockeryYaml(mocks []types.Type) []byte {
	groppedByPackage := make(map[string][]*types.Named)
	for _, m := range mocks {
		named, ok := m.(*types.Named)
		if !ok {
			continue
		}

		_, ok = named.Underlying().(*types.Interface)
		if !ok {
			continue
		}

		pkg := named.Obj().Pkg().Path()
		groppedByPackage[pkg] = append(groppedByPackage[pkg], named)
	}

	// mockery doesn't want work with \t
	buf, emit := createEmmiter()
	emit("with-expecter: True\n")
	emit("packages:\n")

	hasDefines := false
	for pkg, group := range groppedByPackage {
		emit("  %s:\n", pkg)
		emit("    interfaces:\n")

		alreadyDefined := make(map[string]bool)
		for _, i := range group {
			name := i.Obj().Name()
			if !alreadyDefined[name] {
				emit("      %s:\n", name)
				alreadyDefined[name] = true
				hasDefines = true
			}
		}
	}

	if hasDefines {
		return buf.Bytes()
	} else {
		return nil
	}
}

func createMocks(
	interfacesThatNeededMock []types.Type,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) []*GeneratedFunc {
	alreadyCreatedMocks := make(map[types.Type]*GeneratedFunc)
	for _, i := range interfacesThatNeededMock {
		createMock(i, alreadyCreatedMocks, 1, typeContext, qualifier, relativePackagePath, maxDepth)
	}

	var funcs []*GeneratedFunc
	for _, mock := range alreadyCreatedMocks {
		funcs = append(funcs, mock)
	}

	return funcs
}

func createMock(
	typeThatNeededMock types.Type,
	created map[types.Type]*GeneratedFunc,
	depth int,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) {
	switch t := typeThatNeededMock.(type) {
	case *types.Named:
		createNamedMock(t, created, depth, typeContext, qualifier, relativePackagePath, maxDepth)
	case *types.Signature:
		createFuncMock(nil, t, created, depth, typeContext, qualifier, relativePackagePath, maxDepth)
	}
}

func createNamedMock(
	interfacesThatNeededMock *types.Named,
	created map[types.Type]*GeneratedFunc,
	depth int,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) {
	switch u := interfacesThatNeededMock.Underlying().(type) {
	case *types.Interface:
		createInterfaceMock(interfacesThatNeededMock, created, depth, typeContext, qualifier, relativePackagePath, maxDepth)
	case *types.Signature:
		createFuncMock(interfacesThatNeededMock, u, created, depth, typeContext, qualifier, relativePackagePath, maxDepth)
	}
}

func createInterfaceMock(
	interfacesThatNeededMock *types.Named,
	created map[types.Type]*GeneratedFunc,
	depth int,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) {
	iface, ok := interfacesThatNeededMock.Underlying().(*types.Interface)
	if !ok {
		return
	}

	if _, ok := created[interfacesThatNeededMock]; ok {
		return
	}

	obj := interfacesThatNeededMock.Obj()
	if obj.Pkg() == nil {
		return
	}

	// Maybe exists better sln for this case
	var importPath string
	if relativePackagePath == "" {
		importPath = fmt.Sprintf("mocks/%s", obj.Pkg().Path())
	} else {
		importPath = fmt.Sprintf("%s/mocks/%s", relativePackagePath, obj.Pkg().Path())
	}

	importPrefix := qualifier.AddImport("mocks", importPath)

	interfaceName := interfacesThatNeededMock.Obj().Name()
	returnType := types.TypeString(interfacesThatNeededMock, qualifier.Qualifier)
	funcName := fmt.Sprintf("fabric_mock_interface_%s", interfaceName)
	genFunc := &GeneratedFunc{
		Name:       funcName,
		ReturnType: returnType,
	}
	created[interfacesThatNeededMock] = genFunc
	typeContext.AddType(interfacesThatNeededMock)

	if depth == maxDepth {
		emptyMock := createEmptyInterfaceMock(interfacesThatNeededMock, importPrefix, qualifier)
		created[interfacesThatNeededMock] = emptyMock

		return
	}

	var returnValues []*types.Var
	var supportedMethods []*types.Func
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		sig := method.Type().(*types.Signature)

		isSupported := true
		for j := 0; j < sig.Results().Len(); j++ {
			if !isSupported {
				break
			}

			result := sig.Results().At(j)

			namedOrSignatures := getNamedOrSignatureTypes(result.Type())
			for _, t := range namedOrSignatures {
				createMock(t, created, depth+1, typeContext, qualifier, relativePackagePath, maxDepth)
			}

			if !typeContext.IsSupported(result.Type()) {
				isSupported = false
				break
			}
		}

		if isSupported {
			supportedMethods = append(supportedMethods, method)

			for j := 0; j < sig.Results().Len(); j++ {
				result := sig.Results().At(j)
				returnValues = append(returnValues, result)
			}
		}
	}

	buf, emit := createEmmiter()

	varContext := NewVariablesContext(qualifier)
	var names []string

	emit("func %s(\n", funcName)
	emit("\tt *testing.T,\n")
	for _, returnValue := range returnValues {
		name := varContext.CreateUniqueNameForVariable(returnValue)
		names = append(names, name)

		emit("\t%s %s,\n", name, types.TypeString(returnValue.Type(), qualifier.Qualifier))
	}
	emit(") %s {\n", returnType)

	emit("\tgenMock := %s.NewMock%s(t)\n", importPrefix, interfaceName)

	varIndex := 0
	for _, method := range supportedMethods {
		sig := method.Type().(*types.Signature)

		emit("\tgenMock.\n")
		emit("\t\tOn(\"%s\"", method.Name())
		for j := 0; j < sig.Params().Len(); j++ {
			param := sig.Params().At(j)
			paramType := types.TypeString(param.Type(), qualifier.Qualifier)
			emit(", mock.AnythingOfType(\"%s\")", paramType)
		}
		emit(").\n")

		sigString := types.TypeString(sig, qualifier.Qualifier)
		emit("\t\tReturn(%s {\n", sigString)
		emit("\t\t\treturn ")
		for j := 0; j < sig.Results().Len(); j++ {
			if j == 0 {
				emit("%s", names[varIndex])
			} else {
				emit(", %s", names[varIndex])
			}
			varIndex++
		}
		emit("\n\t\t})\n")
	}

	emit("\treturn genMock\n")
	emit("}\n")

	genFunc.Body = buf.String()
}

var Index int = 0

func createFuncMock(
	namedType *types.Named,
	signatureThatNeededMock *types.Signature,
	created map[types.Type]*GeneratedFunc,
	depth int,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) *GeneratedFunc {
	if _, ok := created[signatureThatNeededMock]; ok {
		return nil
	}

	if depth == maxDepth {
		return nil
	}

	sigTypeString := types.TypeString(signatureThatNeededMock, qualifier.Qualifier)

	var returnType string
	if namedType != nil {
		returnType = types.TypeString(namedType, qualifier.Qualifier)
	} else {
		returnType = sigTypeString
	}

	Index++
	funcName := fmt.Sprintf("fabric_mock_func_%d", Index)
	genFunc := &GeneratedFunc{
		Name:       funcName,
		ReturnType: returnType,
	}

	isSupported := true
	for j := 0; j < signatureThatNeededMock.Results().Len(); j++ {
		if !isSupported {
			break
		}

		result := signatureThatNeededMock.Results().At(j)

		namedOrSignatures := getNamedOrSignatureTypes(result.Type())
		for _, t := range namedOrSignatures {
			createMock(t, created, depth+1, typeContext, qualifier, relativePackagePath, maxDepth)
		}

		if !typeContext.IsSupported(result.Type()) {
			isSupported = false
			break
		}
	}

	if !isSupported {
		return nil
	}

	if namedType != nil {
		created[namedType] = genFunc
		typeContext.AddType(namedType)
	} else {
		created[signatureThatNeededMock] = genFunc
		typeContext.AddType(signatureThatNeededMock)
	}

	buf, emit := createEmmiter()

	varContext := NewVariablesContext(qualifier)
	var names []string

	if signatureThatNeededMock.Results().Len() > 1 {
		emit("func %s(\n", funcName)
		for i := 0; i < signatureThatNeededMock.Results().Len(); i++ {
			returnValue := signatureThatNeededMock.Results().At(i)
			name := varContext.CreateUniqueNameForVariable(returnValue)
			names = append(names, name)

			emit("\t%s %s,\n", name, types.TypeString(returnValue.Type(), qualifier.Qualifier))
		}
		emit(") %s {\n", returnType)
	} else {
		emit("func %s(", funcName)
		for i := 0; i < signatureThatNeededMock.Results().Len(); i++ {
			returnValue := signatureThatNeededMock.Results().At(i)
			name := varContext.CreateUniqueNameForVariable(returnValue)
			names = append(names, name)

			emit("%s %s", name, types.TypeString(returnValue.Type(), qualifier.Qualifier))
		}
		emit(") %s {\n", returnType)
	}

	emit("\treturn %s", sigTypeString)

	emit(" {\n")

	emit("\t\treturn ")
	for i := 0; i < signatureThatNeededMock.Results().Len(); i++ {
		if i == 0 {
			emit("%s", names[i])
		} else {
			emit(", %s", names[i])
		}
	}
	emit("\n\t}\n")
	emit("}\n")

	genFunc.Body = buf.String()

	return genFunc
}

func createEmptyInterfaceMock(
	named *types.Named,
	importPrefix string,
	qualifier *ImportQualifier,
) *GeneratedFunc {
	buf, emit := createEmmiter()

	interfaceName := named.Obj().Name()
	returnType := types.TypeString(named, qualifier.Qualifier)
	funcName := fmt.Sprintf("fabric_mock_%s", interfaceName)
	emit("func %s(t *testing.T) %s {\n", funcName, returnType)
	emit("\treturn %s.NewMock%s(t)\n", importPrefix, interfaceName)
	emit("}\n")

	return &GeneratedFunc{
		Name:       funcName,
		Body:       buf.String(),
		ReturnType: returnType,
	}
}
