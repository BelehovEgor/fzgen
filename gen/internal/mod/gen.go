package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
)

const maxDepth = 3

func (s *Struct) Initialize(
	name string,
	defaultQualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
) string {
	var buffer bytes.Buffer
	var w io.Writer = &buffer
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	emit("%s := ", name)
	emit("\t%s\n", CreateConstructor(s, defaultQualifier))

	return buffer.String()
}

func (i Interface) GetConstructors(
	varName string,
	defaultQualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
) []string {
	constructors := make([]string, len(i.Implementations))
	for i, imp := range i.Implementations {
		constructors[i] = imp.Initialize(varName, defaultQualifier, supportedInterfaces, existedFuncs)
	}

	return constructors
}

type GeneratedFunc struct {
	Name, Body  string
	ArgRequired bool
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
	defaultQualifier types.Qualifier,
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
	interfaceTypeString := pkgInterface.TypeString(defaultQualifier)
	if len(pkgInterface.Implementations) == 1 {
		emit("func %s() %s {\n", funcName, interfaceTypeString)
		emit("\treturn %s\n}\n", CreateConstructor(pkgInterface.Implementations[0], defaultQualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, interfaceTypeString)
		emit("\tswitch num %s %d {\n", "%", len(pkgInterface.Implementations))
		for i, s := range pkgInterface.Implementations {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", CreateConstructor(s, defaultQualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	return &GeneratedFunc{
		Name:        funcName,
		Body:        buf.String(),
		ArgRequired: len(pkgInterface.Implementations) > 1,
	}
}

func CreateFabricOfFuncs(
	signature *types.Signature,
	prefix string,
	defaultQualifier types.Qualifier,
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
	typeString := types.TypeString(signature, defaultQualifier)
	if len(suitable) == 1 {
		emit("func %s() %s {\n", funcName, typeString)
		emit("\treturn %s\n}\n", suitable[0].TypeString(defaultQualifier))
	} else {
		emit("func %s(num int) %s {\n", funcName, typeString)
		emit("\tswitch num %s %d {\n", "%", len(suitable))
		for i, s := range suitable {
			emit("\tcase %d:\n", i)
			emit("\t\treturn %s\n", s.TypeString(defaultQualifier))
		}
		emit("\tdefault:\n")
		emit("\t\tpanic(\"unreachable\")\n\t}\n}\n")
	}

	return &GeneratedFunc{
		Name:        funcName,
		Body:        buf.String(),
		ArgRequired: len(suitable) > 1,
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
