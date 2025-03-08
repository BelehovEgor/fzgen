package llm

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"

	"github.com/BelehovEgor/fzgen/gen/internal/mod"
)

var (
	commonRequirements = `
Return new fuzzing target, that process all invalid func call results as error

Requirements:
+ return only fuzz target code
+ without imports
+ no explanation
+ process all edge cases
+ if arguments is invalid, target function shouldn't be call, this case should be skipped without log
+ if there is an explicit exception creation, skip only them, the rest should cause a fuzzing test error
+ situations that should not occur during the execution of the function should end with t.Error
+ don't use not exported fields in validation checks
+ !MISSING in example is problem of encoding, don't use it in result
+ use UTF-8
	`

	patternTargetUseFuzz string = `
You write fuzzing test using golang

Purpose: extend target func result check in fuzzing target 

%s

Target func:

%s

Example of using func:

%s

Current fuzzing target test:

%s

Current import aliases:

%s
	`

	patternTargetFuzz string = `
You write fuzzing test using golang

Purpose: extend target func result check in fuzzing target 

%s

Target func:

%s

Current fuzzing target test:

%s

Current import aliases:

%s
	`

	patternTarget string = `
You write fuzzing test using golang

%s

Purpose: generate fuzz target func for function:

%s

%s
	`
)

func CreatePrompt(
	fset *token.FileSet,
	target *mod.Func,
	fuzzTarget string,
	qualifier *mod.ImportQualifier,
) string {
	if len(fuzzTarget) == 0 {
		return fillTemplate(patternTarget, getSourceCode(target.AstFuncDecl, fset))
	}

	if len(target.Uses) > 0 {
		return fillTemplate(
			patternTargetUseFuzz,
			getSourceCode(target.AstFuncDecl, fset),
			getSourceCode(target.Uses[0], fset),
			fuzzTarget,
			getImports(qualifier),
		)
	}

	return fillTemplate(patternTargetFuzz, getSourceCode(target.AstFuncDecl, fset), fuzzTarget, getImports(qualifier))
}

func getSourceCode(target *ast.FuncDecl, fset *token.FileSet) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, target); err != nil {
		log.Fatalf("Failed to print function declaration: %v", err)
	}
	return buf.String()
}

func getImports(qualifier *mod.ImportQualifier) string {
	buf, emit := mod.CreateEmmiter()
	for _, str := range qualifier.GetImportStrings() {
		emit("%s\n", str)
	}

	return buf.String()
}

func fillTemplate(template string, a ...any) string {
	return fmt.Sprintf(template, commonRequirements, a)
}
