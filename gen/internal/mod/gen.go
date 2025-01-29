package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"sort"
)

const maxDepth = 10

type GeneratedFunc struct {
	Name, Body, ReturnType string
}

func GenerateFabrics(
	targets []*Func,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
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

	for i, target := range targets {
		for j, param := range target.Params() {
			switch u := param.Type().Underlying().(type) {
			case *types.Interface:
				if funcs := createFabricOfInterfaces(param, qualifier, supportedInterfaces); funcs != nil {
					generated[types.TypeString(u, qualifier.Qualifier)] = funcs
					typeContext.AddType(param.Type())
				}
			case *types.Struct:
				for notNativeType := range getNotNativeTypes(u, 1) {
					switch t := notNativeType.Type().Underlying().(type) {
					case *types.Interface:
						if funcs := createFabricOfInterfaces(
							notNativeType,
							qualifier,
							supportedInterfaces,
						); funcs != nil {
							generated[types.TypeString(t, qualifier.Qualifier)] = funcs
							typeContext.AddType(t)
						}
					case *types.Signature:
						if f := createFabricOfFuncs(
							t,
							fmt.Sprintf("%d_%d", i, j),
							qualifier,
							existedFuncs,
						); f != nil {
							generated[types.TypeString(t, qualifier.Qualifier)] = append(generated[types.TypeString(t, qualifier.Qualifier)], f)
							typeContext.AddType(t)
						}
					}
				}
			case *types.Signature:
				if f := createFabricOfFuncs(
					u,
					fmt.Sprintf("%d_%d", i, j),
					qualifier,
					existedFuncs,
				); f != nil {
					generated[types.TypeString(u, qualifier.Qualifier)] = append(generated[types.TypeString(u, qualifier.Qualifier)], f)

					typeContext.AddType(param.Type())
				}
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
			return funcs[i].Name > funcs[j].Name
		})

		for _, f := range funcs {
			emit(
				"FabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
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
			"FabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
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

func createFabricOfEmptyInterface(
	existedStructs map[string]*Struct,
	qualifier *ImportQualifier,
) []*GeneratedFunc {
	if len(existedStructs) == 0 {
		return nil
	}

	varContext := NewVariablesContext(qualifier)

	var structs []*Struct
	for _, s := range existedStructs {
		structs = append(structs, s)
	}

	sort.Slice(structs, func(i, j int) bool {
		return structs[i].StructName > structs[j].StructName
	})

	funcs := make([]*GeneratedFunc, len(structs))

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

	return funcs
}

func createFabricOfInterfaces(
	v *types.Var,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
) []*GeneratedFunc {
	typeString := v.Type().String()
	supported, ok := supportedInterfaces[typeString]
	if !ok {
		return nil
	}
	return createFabricOfInterfaces2(supported, qualifier)
}

func createFabricOfInterfaces2(
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
		}
	}

	return funcs
}

func createFabricOfFuncs(
	signature *types.Signature,
	prefix string,
	qualifier *ImportQualifier,
	existedFuncs []*Func,
) *GeneratedFunc {
	suitable := FindSuitables(signature, existedFuncs)
	if len(suitable) == 0 {
		return nil
	}

	sort.Slice(suitable, func(i, j int) bool {
		return suitable[i].FuncName > suitable[j].FuncName
	})

	buf, emit := createEmmiter()

	funcName := "fabric_func_" + prefix
	typeString := types.TypeString(signature, qualifier.Qualifier)
	if len(suitable) == 1 {
		emit("func %s() %s {\n", funcName, typeString)
		emit("\treturn %s\n}\n", suitable[0].TypeString(qualifier.Qualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, typeString)
		emit("\tswitch num %s %d {\n", "%", len(suitable))
		for i, s := range suitable {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", s.TypeString(qualifier.Qualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	return &GeneratedFunc{
		Name:       funcName,
		Body:       buf.String(),
		ReturnType: types.TypeString(signature, qualifier.Qualifier),
	}
}
