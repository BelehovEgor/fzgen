package mod

import (
	"fmt"
	"go/types"
	"strings"
)

type VariablesContext struct {
	importQualifier *ImportQualifier

	varNames map[string]int

	uniqueNumber int
	Fabrics      map[types.Type]*GeneratedFunc
}

func NewVariablesContext(importQualifier *ImportQualifier) *VariablesContext {
	varNames := make(map[string]int)
	// all key word
	varNames["package"] = 0

	varNames["data"] = 0 // data []bytes
	varNames["t"] = 0    // t testing.T

	return &VariablesContext{
		importQualifier: importQualifier,
		uniqueNumber:    1,
		varNames:        varNames,
	}
}

func (vc *VariablesContext) CreateUniqueNameForVariable(v *types.Var) string {
	wanted := v.Name()
	if wanted == "" {
		parts := strings.Split(v.Type().String(), ".")

		lastPart := strings.ToLower(parts[len(parts)-1])
		lastPart = strings.Replace(lastPart, "{}", "", -1)
		lastPart = strings.Replace(lastPart, "[]", "_arr_", -1)

		wanted = lastPart + "_"
	}

	return vc.CreateUniqueName(wanted)
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
