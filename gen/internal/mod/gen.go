package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"sort"
)

type GeneratedFunc struct {
	Name, Body, ReturnType string
	Type                   types.Type
}

func GenerateFabrics(
	targets []*Func,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
	maxDepth int,
) map[string][]*GeneratedFunc {
	generated := make(map[string][]*GeneratedFunc)

	supportedInterfaces := make(map[string]*Interface)
	for i := range typeContext.SupportedInterfaces {
		supportedInterfaces[i.TypesNamed.String()] = i
	}

	existedStructs := make(map[string]*Struct)
	for _, s := range typeContext.ExistedStructs {
		existedStructs[s.TypesNamed.String()] = s
	}

	existedFuncs := typeContext.ExistedFuncs

	emptyInterfaceFabrics := createFabricOfEmptyInterface(existedStructs, qualifier)
	generated["interface{}"] = emptyInterfaceFabrics

	for _, target := range targets {
		for _, param := range target.Params() {
			if funcs := createFabrics(param, qualifier, supportedInterfaces, existedFuncs, maxDepth); funcs != nil {
				generated[types.TypeString(param.Type(), qualifier.Qualifier)] = funcs
				typeContext.AddType(param.Type())
			}
		}
	}

	return generated
}

func GenerateInitTestFunc(
	generatedFuncs map[string][]*GeneratedFunc,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
) *GeneratedFunc {
	buf, emit := createEmmiter()

	emit("var FabricFuncsForCustomTypes map[string][]reflect.Value\n\n")
	emit("func TestMain(m *testing.M) {\n")
	emit("\tFabricFuncsForCustomTypes = make(map[string][]reflect.Value)\n")

	for _, funcs := range generatedFuncs {
		sort.Slice(funcs, func(i, j int) bool {
			return funcs[i].ReturnType > funcs[j].ReturnType &&
				funcs[i].Name > funcs[j].Name
		})

		for _, f := range funcs {
			emit(
				"\tFabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
				f.ReturnType,
				f.ReturnType,
				f.Name,
			)
		}
	}

	var orderedConstructors []*Constructor
	for c := range typeContext.ValidConstructors {
		orderedConstructors = append(orderedConstructors, c)
	}

	sort.Slice(orderedConstructors, func(i, j int) bool {
		return orderedConstructors[i].Func.FuncName > orderedConstructors[j].Func.FuncName
	})

	for _, constructor := range orderedConstructors {
		if !constructor.Func.TypesFunc.Exported() && !qualifier.isLocalTest {
			continue
		}

		returnType := types.TypeString(constructor.ReturnType, qualifier.Qualifier)

		emit(
			"\tFabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
			returnType,
			returnType,
			constructor.Func.TypeString(qualifier.Qualifier),
		)
	}

	emit("\tm.Run()\n")
	emit("}")

	return &GeneratedFunc{
		Name: "TestMain",
		Body: buf.String(),
	}
}

func createEmmiter() (*bytes.Buffer, func(format string, args ...interface{})) {
	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	return buf, emit
}

func createFabrics(
	param *types.Var,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
	maxDepth int,
) []*GeneratedFunc {
	if maxDepth == 0 {
		return nil
	}

	// get types that required fabrics
	baseTypes := getNamedOrSignatureTypes(param.Type())
	if len(baseTypes) == 0 {
		return nil
	}

	var fabrics []*GeneratedFunc
	for _, base := range baseTypes {
		baseFabrics := createTypeFabrics(base, qualifier, supportedInterfaces, existedFuncs, maxDepth)
		if baseFabrics != nil {
			fabrics = append(fabrics, baseFabrics...)
		}
	}

	return fabrics
}

func createTypeFabrics(
	paramType types.Type,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
	maxDepth int,
) []*GeneratedFunc {
	switch t := paramType.(type) {
	case *types.Named:
		return createNamedFabrics(t, qualifier, supportedInterfaces, existedFuncs, maxDepth)
	case *types.Signature:
		return createFabricOfFuncs(t, nil, qualifier, existedFuncs)
	}

	return nil
}

func createNamedFabrics(
	named *types.Named,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
	maxDepth int,
) []*GeneratedFunc {
	switch u := named.Underlying().(type) {
	case *types.Interface:
		return createFabricOfInterfaces(named, qualifier, supportedInterfaces)
	case *types.Struct:
		baseTypes := getNotNativeTypes(u, 0)
		if len(baseTypes) == 0 {
			return nil
		}

		var result []*GeneratedFunc
		for base := range baseTypes {
			fieldFabrics := createTypeFabrics(base, qualifier, supportedInterfaces, existedFuncs, maxDepth-1)
			result = append(result, fieldFabrics...)
		}

		return result
	case *types.Signature:
		return createFabricOfFuncs(u, named, qualifier, existedFuncs)
	}

	return nil
}

func createFabricOfEmptyInterface(
	existedStructs map[string]*Struct,
	qualifier *ImportQualifier,
) []*GeneratedFunc {
	varContext := NewVariablesContext(qualifier)

	var structs []*Struct
	for _, s := range existedStructs {
		structs = append(structs, s)
	}

	sort.Slice(structs, func(i, j int) bool {
		return structs[i].StructName > structs[j].StructName
	})

	funcs := make([]*GeneratedFunc, len(structs)+1)

	for i, impl := range structs {
		buf, emit := createEmmiter()
		typeName := varContext.CreateUniqueName(impl.StructName)
		funcName := fmt.Sprintf("fabric_interface_empty_%s", typeName)
		typeString := impl.TypeString(qualifier.Qualifier)
		emit("func %s(impl %s) interface{} {\n", funcName, typeString)
		emit("\treturn impl\n")
		emit("}\n")

		funcs[i] = &GeneratedFunc{
			Name:       funcName,
			Body:       buf.String(),
			ReturnType: "interface {}",
		}
	}

	buf, emit := createEmmiter()
	funcName := "fabric_interface_empty_string"
	emit("func %s(impl string) interface{} {\n", funcName)
	emit("\treturn impl\n")
	emit("}\n")

	funcs[len(structs)] = &GeneratedFunc{
		Name:       funcName,
		Body:       buf.String(),
		ReturnType: "interface {}",
	}

	return funcs
}

func createFabricOfInterfaces(
	named *types.Named,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
) []*GeneratedFunc {
	typeString := named.String()
	supported, ok := supportedInterfaces[typeString]
	if !ok {
		return nil
	}
	return createFabricOfInterfaces2(named, supported, qualifier)
}

func createFabricOfInterfaces2(
	named *types.Named,
	pkgInterface *Interface,
	qualifier *ImportQualifier,
) []*GeneratedFunc {
	if len(pkgInterface.Implementations) == 0 {
		return nil
	}

	varContext := NewVariablesContext(qualifier)
	interfaceTypeString := pkgInterface.TypeString(qualifier.Qualifier)

	sort.Slice(pkgInterface.Implementations, func(i, j int) bool {
		return pkgInterface.Implementations[i].StructName > pkgInterface.Implementations[j].StructName
	})

	funcs := make([]*GeneratedFunc, len(pkgInterface.Implementations))

	for i, impl := range pkgInterface.Implementations {
		buf, emit := createEmmiter()
		funcSuffix := varContext.CreateUniqueName(impl.StructName)
		funcName := fmt.Sprintf("fabric_interface_%s_%s_%s", pkgInterface.PkgName, pkgInterface.InterfaceName, funcSuffix)
		emit("func %s(impl ", funcName)
		if impl.AsPointer {
			emit("*")
		}
		typeString := impl.TypeString(qualifier.Qualifier)
		emit("%s) %s {\n", typeString, interfaceTypeString)
		emit("\treturn impl\n")
		emit("}\n")

		funcs[i] = &GeneratedFunc{
			Name:       funcName,
			Body:       buf.String(),
			ReturnType: pkgInterface.TypeString(qualifier.Qualifier),
			Type:       named,
		}
	}

	return funcs
}

func createFabricOfFuncs(
	signature *types.Signature,
	named *types.Named,
	qualifier *ImportQualifier,
	existedFuncs []*Func,
) []*GeneratedFunc {
	suitable := FindSuitables(signature, existedFuncs)
	if len(suitable) == 0 {
		return nil
	}

	sort.Slice(suitable, func(i, j int) bool {
		return suitable[i].FuncName > suitable[j].FuncName
	})

	buf, emit := createEmmiter()

	Index++
	funcName := fmt.Sprintf("fabric_func_%d", Index)
	returnType := types.TypeString(signature, qualifier.Qualifier)
	if named != nil {
		returnType = types.TypeString(named, qualifier.Qualifier)
	}

	if len(suitable) == 1 {
		emit("func %s() %s {\n", funcName, returnType)
		emit("\treturn %s\n}\n", suitable[0].TypeString(qualifier.Qualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, returnType)
		emit("\tswitch num %s %d {\n", "%", len(suitable))
		for i, s := range suitable {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", s.TypeString(qualifier.Qualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	result := make([]*GeneratedFunc, 1)
	result[0] = &GeneratedFunc{
		Name:       funcName,
		Body:       buf.String(),
		ReturnType: returnType,
		Type:       named,
	}
	return result
}
