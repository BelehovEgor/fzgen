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
	var interfacesThatNeededMock []*types.Named
	var funcsThatNeededMock []*types.Signature

	for _, target := range targets {
		for _, param := range target.Params() {
			baseTypes := getNamedOrSignatureTypes(param.Type())
			for _, baseType := range baseTypes {
				switch t := baseType.(type) {
				case *types.Named:
					switch u := t.Underlying().(type) {
					case *types.Interface:
						interfacesThatNeededMock = append(interfacesThatNeededMock, t)
					case *types.Struct:
						for notNativeType := range getNotNativeTypes(u, 1) {
							switch structNotNative := notNativeType.(type) {
							case *types.Named:
								if types.IsInterface(structNotNative) {
									interfacesThatNeededMock = append(interfacesThatNeededMock, structNotNative)
								}
							case *types.Signature:
								funcsThatNeededMock = append(funcsThatNeededMock, structNotNative)
							}
						}
					}
				case *types.Signature:
					funcsThatNeededMock = append(funcsThatNeededMock, t)
				}
			}
		}
	}

	interfaceMockFabrics := createInterfaceMocks(interfacesThatNeededMock, typeContext, qualifier, relativePackagePath, maxDepth)

	return interfaceMockFabrics, generateMockeryYaml(interfacesThatNeededMock)
}

func generateMockeryYaml(interfaces []*types.Named) []byte {
	groppedByPackage := make(map[string][]*types.Named)
	for _, i := range interfaces {
		pkg := i.Obj().Pkg().Path()
		groppedByPackage[pkg] = append(groppedByPackage[pkg], i)
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

func createInterfaceMocks(
	interfacesThatNeededMock []*types.Named,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	relativePackagePath string,
	maxDepth int,
) []*GeneratedFunc {
	alreadyCreatedMocks := make(map[*types.Named]*GeneratedFunc)
	for _, i := range interfacesThatNeededMock {
		createInterfaceMockRec(i, alreadyCreatedMocks, 1, typeContext, qualifier, relativePackagePath, maxDepth)
	}

	var funcs []*GeneratedFunc
	for _, mock := range alreadyCreatedMocks {
		funcs = append(funcs, mock)
	}

	return funcs
}

func createInterfaceMockRec(
	interfacesThatNeededMock *types.Named,
	created map[*types.Named]*GeneratedFunc,
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
	funcName := fmt.Sprintf("fabric_mock_%s", interfaceName)
	genFunc := &GeneratedFunc{
		Name:       funcName,
		ReturnType: returnType,
	}
	created[interfacesThatNeededMock] = genFunc
	typeContext.AddType(interfacesThatNeededMock)

	if depth == maxDepth {
		emptyMock := createEmptyMock(interfacesThatNeededMock, importPrefix, qualifier)
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
				named, ok := t.(*types.Named)
				if !ok {
					continue
				}

				createInterfaceMockRec(named, created, depth+1, typeContext, qualifier, relativePackagePath, maxDepth)
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

func createEmptyMock(
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
