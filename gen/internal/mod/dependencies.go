package mod

import (
	"fmt"
	"go/types"
)

func (s *Struct) GetNotNativeTypes() map[*types.Var]bool {
	return getNotNativeTypes(s.TypesStruct, 1)
}

func HasNotNative(s *Struct) bool {
	return hasNotNative(s.TypesStruct, 1)
}

func GetInputParams(function *Func) []*types.Var {
	f := function.TypesFunc
	wrappedSig := f.Type().(*types.Signature)
	var inputParams []*types.Var

	for i := 0; i < wrappedSig.Params().Len(); i++ {
		v := wrappedSig.Params().At(i)
		inputParams = append(inputParams, v)
	}

	return inputParams
}

func getNotNativeTypes(s *types.Struct, depth int) map[*types.Var]bool {
	if depth > maxDepth {
		return nil
	}

	notNative := make(map[*types.Var]bool)
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		switch u := field.Type().Underlying().(type) {
		case *types.Interface, *types.Signature:
			notNative[field] = true
		case *types.Struct:
			fieldNotNative := getNotNativeTypes(u, depth+1)
			if fieldNotNative == nil {
				continue
			}

			for v := range fieldNotNative {
				notNative[v] = true
			}
		}
	}

	return notNative
}

func hasNotNative(s *types.Struct, depth int) bool {
	if depth > maxDepth {
		return false
	}

	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		switch field.Type().Underlying().(type) {
		case *types.Interface, *types.Signature:
			return true
		case *types.Struct:
			has := hasNotNative(s, depth+1)
			if has {
				return true
			}
		}
	}

	return false
}

type ImportQualifier struct {
	pkgName, pkgPath, outPkgPath, outPkgName string
	isLocalTest                              bool

	Imports     map[string]string
	importNames map[string]int
}

func CreateQualifier(pkgName, pkgPath, outPkgName, outPkgPath string, isLocalTest bool) *ImportQualifier {
	q := &ImportQualifier{
		pkgName:     pkgName,
		pkgPath:     pkgPath,
		outPkgPath:  outPkgPath,
		outPkgName:  outPkgName,
		isLocalTest: isLocalTest,
		Imports:     make(map[string]string),
		importNames: make(map[string]int),
	}

	q.Imports[outPkgPath] = ""
	q.importNames[outPkgName] = 0

	if isLocalTest {
		q.Imports[pkgPath] = ""
		q.importNames[pkgName] = 0
	} else {
		if pkgName == outPkgName {
			q.Imports[pkgPath] = fmt.Sprintf("%s_1", pkgName)
			q.importNames[pkgName] = 1
		} else {
			q.Imports[pkgPath] = pkgName
			q.importNames[pkgName] = 0
		}
	}

	if pkgName != "fuzzer" {
		q.importNames[pkgName] = 0
		q.Imports["github.com/thepudds/fzgen/fuzzer"] = ""
	} else {
		q.importNames[pkgName] = 1
		q.Imports["github.com/thepudds/fzgen/fuzzer"] = "fuzzer_1"
	}

	return q
}

func (iq *ImportQualifier) GetImportStrings() []string {
	var imports []string

	imports = append(imports, "\"testing\"")

	for path, name := range iq.Imports {
		if path == iq.outPkgPath {
			continue
		}

		if name == "" {
			imports = append(imports, fmt.Sprintf("\"%s\"", path))
		} else {
			imports = append(imports, fmt.Sprintf("%s \"%s\"", name, path))
		}
	}

	return imports
}

func (iq *ImportQualifier) Qualifier(p *types.Package) string {
	defaultName := p.Name()
	path := p.Path()

	if iq.isLocalTest && path == iq.pkgPath {
		return ""
	}

	name, has := iq.Imports[path]
	if has {
		if name == "" {
			return defaultName
		}
		return name
	}

	idx, has := iq.importNames[defaultName]
	if has {
		idx++
		iq.importNames[defaultName] = idx

		name = fmt.Sprintf("%s_%d", defaultName, idx)
		iq.Imports[path] = name

		return name
	}

	iq.Imports[path] = ""
	iq.importNames[defaultName] = 0

	return defaultName
}

type VariablesContext struct {
	importQualifier *ImportQualifier

	varNames map[string]int

	uniqueNumber int
	Fabrics      map[types.Type]*GeneratedFunc
}

func CreateVariablesContext(importQualifier *ImportQualifier) *VariablesContext {
	varNames := make(map[string]int)
	// all key word
	varNames["package"] = 0

	return &VariablesContext{
		importQualifier: importQualifier,
		uniqueNumber:    1,
		varNames:        varNames,
	}
}

func (vc *VariablesContext) CreateUniqueName(wanted string) string {
	if idx, ok := vc.varNames[wanted]; ok {
		vc.varNames[wanted] = idx + 1
		return fmt.Sprintf("%s_%d", wanted, idx)
	}

	if idx, ok := vc.importQualifier.importNames[wanted]; ok {
		vc.varNames[wanted] = idx + 1
		return fmt.Sprintf("%s_%d", wanted, idx)
	}

	vc.varNames[wanted] = 1
	return wanted
}
