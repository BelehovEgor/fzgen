package mod

import (
	"fmt"
	"go/types"
)

type ImportQualifier struct {
	pkgName, pkgPath, outPkgPath, outPkgName string
	isLocalTest                              bool

	Imports     map[string]string
	importNames map[string]int
}

func NewQualifier(pkgName, pkgPath, outPkgName, outPkgPath string, isLocalTest bool) *ImportQualifier {
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
		q.Imports["github.com/BelehovEgor/fzgen/fuzzer"] = ""
	} else {
		q.importNames[pkgName] = 1
		q.Imports["github.com/BelehovEgor/fzgen/fuzzer"] = "fuzzer_1"
	}

	q.Imports["reflect"] = ""
	q.Imports["github.com/stretchr/testify/mock"] = ""

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

func (iq *ImportQualifier) AddImport(defaultName, path string) string {
	name, has := iq.Imports[path]
	if has {
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

	iq.Imports[path] = defaultName
	iq.importNames[defaultName] = 0

	return defaultName
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
