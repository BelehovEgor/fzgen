package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
	"sort"
	"strings"
)

const maxDepth = 10

type GeneratedFunc struct {
	Name, Body, ReturnType string
}

func GenerateFabrics(
	targets []*Func,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
) map[string]*GeneratedFunc {
	generated := make(map[string]*GeneratedFunc)

	supportedInterfaces := make(map[string]*Interface)
	for i := range typeContext.SupportedInterfaces {
		supportedInterfaces[i.TypesNamed.String()] = i
	}

	existedStructs := make(map[string]*Struct)
	for _, s := range typeContext.ExistedStructs {
		existedStructs[s.TypesNamed.String()] = s
	}

	existedFuncs := typeContext.ExistedFuncs

	emptyInterfaceFabric := createFabricOfEmptyIntarface(existedStructs, qualifier)
	if emptyInterfaceFabric != nil {
		generated["interface{}"] = emptyInterfaceFabric
	}

	for i, target := range targets {
		for j, param := range target.Params() {
			switch u := param.Type().Underlying().(type) {
			case *types.Interface:
				if f := createFabricOfInterfaces(param, qualifier, supportedInterfaces); f != nil {
					generated[types.TypeString(u, qualifier.Qualifier)] = f
					typeContext.AddType(param.Type())
				}
			case *types.Struct:
				for notNativeType := range getNotNativeTypes(u, 1) {
					switch t := notNativeType.Type().Underlying().(type) {
					case *types.Interface:
						if f := createFabricOfInterfaces(
							notNativeType,
							qualifier,
							supportedInterfaces,
						); f != nil {
							generated[types.TypeString(t, qualifier.Qualifier)] = f
							typeContext.AddType(t)
						}
					case *types.Signature:
						if f := createFabricOfFuncs(
							t,
							fmt.Sprintf("%d_%d", i, j),
							qualifier,
							existedFuncs,
						); f != nil {
							generated[types.TypeString(t, qualifier.Qualifier)] = f
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
					generated[types.TypeString(u, qualifier.Qualifier)] = f
					typeContext.AddType(param.Type())
				}
			}
		}
	}

	return generated
}

func GenerateInitTestFunc(
	generatedFuncs map[string]*GeneratedFunc,
	typeContext *TypeContext,
	qualifier *ImportQualifier,
) *GeneratedFunc {
	buf, emit := createEmmiter()

	emit("var FabricFuncsForCustomTypes map[string][]reflect.Value\n\n")
	emit("func TestMain(m *testing.M) {\n")
	emit("\tFabricFuncsForCustomTypes = make(map[string][]reflect.Value)\n")

	for _, f := range generatedFuncs {
		emit(
			"FabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
			f.ReturnType,
			f.ReturnType,
			f.Name,
		)
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

		sig := constructor.Func.GetSignature()
		returnType := types.TypeString(sig.Results().At(0).Type(), qualifier.Qualifier)

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

func createFabricOfEmptyIntarface(
	existedStructs map[string]*Struct,
	qualifier *ImportQualifier,
) *GeneratedFunc {
	if len(existedStructs) == 0 {
		return nil
	}

	buf, emit := createEmmiter()
	varContext := NewVariablesContext(qualifier)

	emit("func fabric_interface_empty(")

	var structs []*Struct
	for _, s := range existedStructs {
		structs = append(structs, s)
	}

	sort.Slice(structs, func(i, j int) bool {
		return structs[i].StructName > structs[j].StructName
	})

	if len(structs) > 1 {
		emit("\n\tnum int,\n")
		var argNames []string
		for _, s := range structs {
			name := varContext.CreateUniqueName(strings.ToLower(s.StructName))
			argNames = append(argNames, name)
			typeString := s.TypeString(qualifier.Qualifier)
			emit("\t%s %s,\n", name, typeString)
		}
		emit(") interface{} {\n")

		emit("\tswitch num %s %d {\n", "%", len(structs))
		for i, s := range argNames {
			emit("\t\tcase %d:\n", i)
			emit("\t\t\treturn %s\n", s)
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	} else {
		var structName string
		for _, value := range structs {
			typeString := value.TypeString(qualifier.Qualifier)
			structName = strings.ToLower(value.StructName)
			emit("%s %s) interface{} {\n", structName, typeString)

		}

		emit("\treturn %s", structName)
		emit("}\n")
	}

	return &GeneratedFunc{
		Name:       "fabric_interface_empty",
		Body:       buf.String(),
		ReturnType: "interface {}",
	}
}

func createFabricOfInterfaces(
	v *types.Var,
	qualifier *ImportQualifier,
	supportedInterfaces map[string]*Interface,
) *GeneratedFunc {
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
) *GeneratedFunc {
	if len(pkgInterface.Implementations) == 0 {
		return nil
	}

	buf, emit := createEmmiter()

	varContext := NewVariablesContext(qualifier)
	funcName := "fabric_interface_" + pkgInterface.PkgName + pkgInterface.InterfaceName
	interfaceTypeString := pkgInterface.TypeString(qualifier.Qualifier)

	sort.Slice(pkgInterface.Implementations, func(i, j int) bool {
		return pkgInterface.Implementations[i].StructName > pkgInterface.Implementations[j].StructName
	})

	if len(pkgInterface.Implementations) > 1 {
		var argNames []string
		emit("func %s(\n\tnum int,\n", funcName)
		for _, s := range pkgInterface.Implementations {
			name := varContext.CreateUniqueName("impl")
			argNames = append(argNames, name)
			typeString := s.TypeString(qualifier.Qualifier)
			if s.AsPointer {
				emit("\t%s *%s,\n", name, typeString)
			} else {
				emit("\t%s %s,\n", name, typeString)
			}
		}

		emit(") %s {\n", interfaceTypeString)

		emit("\tswitch num %s %d {\n", "%", len(pkgInterface.Implementations))
		for i, s := range argNames {
			emit("\t\tcase %d:\n", i)
			emit("\t\t\treturn %s\n", s)
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	} else {
		emit("func %s(impl ", funcName)
		typeString := pkgInterface.Implementations[0].TypeString(qualifier.Qualifier)
		if pkgInterface.Implementations[0].AsPointer {
			emit("*%s) ", typeString)
		} else {
			emit("%s) ", typeString)
		}
		emit("%s {\n", interfaceTypeString)
		emit("\treturn impl")
		emit("}\n")
	}

	return &GeneratedFunc{
		Name:       funcName,
		Body:       buf.String(),
		ReturnType: pkgInterface.TypeString(qualifier.Qualifier),
	}
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
