package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
)

const maxDepth = 10

type GeneratedFunc struct {
	Name, Body, ReturnType string
	ArgRequired            bool
}

func CreateConstructor(s *Struct, defaultQualifier types.Qualifier) string {
	typeString := s.TypeString(defaultQualifier)

	if s.AsPointer {
		return fmt.Sprintf("&%s{}", typeString)
	}

	return fmt.Sprintf("%s{}", typeString)
}

func CreateFabricOfInterfaces(
	pkgInterface *Interface,
	qualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
) *GeneratedFunc {
	if len(pkgInterface.Implementations) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	funcName := "fabric_interface_" + pkgInterface.PkgName + pkgInterface.InterfaceName
	interfaceTypeString := pkgInterface.TypeString(qualifier)
	if len(pkgInterface.Implementations) == 1 {
		emit("func %s() %s {\n", funcName, interfaceTypeString)
		emit("\treturn %s\n}\n", CreateConstructor(pkgInterface.Implementations[0], qualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, interfaceTypeString)
		emit("\tswitch num %s %d {\n", "%", len(pkgInterface.Implementations))
		for i, s := range pkgInterface.Implementations {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", CreateConstructor(s, qualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	return &GeneratedFunc{
		Name:        funcName,
		Body:        buf.String(),
		ArgRequired: len(pkgInterface.Implementations) > 1,
		ReturnType:  pkgInterface.TypeString(qualifier),
	}
}

func CreateFabricOfFuncs(
	signature *types.Signature,
	prefix string,
	qualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
) *GeneratedFunc {
	suitable := FindSuitables(signature, existedFuncs)
	if len(suitable) == 0 {
		return nil
	}

	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	funcName := "fabric_func_" + prefix
	typeString := types.TypeString(signature, qualifier)
	if len(suitable) == 1 {
		emit("func %s() %s {\n", funcName, typeString)
		emit("\treturn %s\n}\n", suitable[0].TypeString(qualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, typeString)
		emit("\tswitch num %s %d {\n", "%", len(suitable))
		for i, s := range suitable {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", s.TypeString(qualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	return &GeneratedFunc{
		Name:        funcName,
		Body:        buf.String(),
		ArgRequired: len(suitable) > 1,
		ReturnType:  types.TypeString(signature, qualifier),
	}
}

func GenerateFabrics(
	targets []*Func,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
	defaultQualifier types.Qualifier,
) map[types.Type]*GeneratedFunc {
	generated := make(map[types.Type]*GeneratedFunc)
	for i, target := range targets {
		for j, param := range GetInputParams(target) {
			switch u := param.Type().Underlying().(type) {
			case *types.Interface:
				if f := createFabricOfInterfaces(param, defaultQualifier, supportedInterfaces, existedFuncs); f != nil {
					generated[u] = f
				}
			case *types.Struct:
				for notNativeType := range getNotNativeTypes(u, 1) {
					switch t := notNativeType.Type().Underlying().(type) {
					case *types.Interface:
						if f := createFabricOfInterfaces(
							notNativeType,
							defaultQualifier,
							supportedInterfaces,
							existedFuncs,
						); f != nil {
							generated[t] = f
						}
					case *types.Signature:
						if f := CreateFabricOfFuncs(
							t,
							fmt.Sprintf("%d_%d", i, j),
							defaultQualifier,
							supportedInterfaces,
							existedFuncs,
						); f != nil {
							generated[t] = f
						}
					}
				}
			case *types.Signature:
				if f := CreateFabricOfFuncs(
					u,
					fmt.Sprintf("%d_%d", i, j),
					defaultQualifier,
					supportedInterfaces,
					existedFuncs,
				); f != nil {
					generated[u] = f
				}
			}
		}
	}

	return generated
}

func createFabricOfInterfaces(
	v *types.Var,
	defaultQualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
) *GeneratedFunc {
	typeString := v.Type().String()
	supported, ok := supportedInterfaces[typeString]
	if !ok {
		return nil
	}
	return CreateFabricOfInterfaces(supported, defaultQualifier, supportedInterfaces, existedFuncs)
}

func GenerateInitTestFunc(
	generatedFuncs map[types.Type]*GeneratedFunc,
	constructors []*Func,
	qualifier types.Qualifier,
) *GeneratedFunc {
	buf := new(bytes.Buffer)
	var w io.Writer = buf
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

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

	for _, constructor := range constructors {
		sig := constructor.GetSignature()
		if sig.Results().Len() != 1 || sig.Recv() != nil {
			continue
		}

		returnType := types.TypeString(sig.Results().At(0).Type(), qualifier)

		emit(
			"FabricFuncsForCustomTypes[\"%s\"] = append(FabricFuncsForCustomTypes[\"%s\"], reflect.ValueOf(%s))\n",
			returnType,
			returnType,
			constructor.TypeString(qualifier),
		)
	}

	emit("\tm.Run()\n")
	emit("}")

	return &GeneratedFunc{
		Name:        "TestMain",
		Body:        buf.String(),
		ArgRequired: false,
	}
}
