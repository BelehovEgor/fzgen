package mod

import (
	"bytes"
	"fmt"
	"go/types"
	"io"
)

const maxDepth = 3

func (s Struct) Initialize(
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
	emit("\t%s\n", s.init(defaultQualifier, supportedInterfaces, existedFuncs, 1))

	return buffer.String()
}

func (s Struct) init(
	defaultQualifier types.Qualifier,
	supportedInterfaces map[string]*Interface,
	existedFuncs []*Func,
	depth int,
) string {
	var buffer bytes.Buffer
	var w io.Writer = &buffer
	emit := func(format string, args ...interface{}) {
		fmt.Fprintf(w, format, args...)
	}

	if depth > maxDepth {
		return ""
	}

	if s.AsPointer {
		emit("&")
	}

	emit("%s{", s.TypeString(defaultQualifier))
	notNative := s.GetNotNativeTypes()
	for _, v := range notNative {
		switch u := v.Type().Underlying().(type) {
		case *types.Interface:
			typeName := v.Type().String()
			supported, ok := supportedInterfaces[typeName]
			if !ok {
				continue
			}
			init := supported.Implementations[0].init(defaultQualifier, supportedInterfaces, existedFuncs, depth+1)
			if init != "" {
				emit("\n%s: %s,", v.Name(), init)
			}

		case *types.Signature:
			suitable := FindSuitables(u, existedFuncs)
			if len(suitable) == 0 {
				continue
			}
			emit(" %s: %s,", v.Name(), suitable[0].TypeString(defaultQualifier))
		}
	}
	emit("\n}")

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
