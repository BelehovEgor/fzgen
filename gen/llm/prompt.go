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
	patternTargetUseFuzz string = `
You write fuzzing test using golang

Purpose: extend target func result check in fuzzing target 

Target func:

%s

Example of using func:

%s

Current fuzzing target test:

%s

Current import aliases:

%s

Return new fuzzing target, that process all available func call results 

Requirements:
+ only code
+ no explanation
+ process all edge cases
+ if arguments is invalid, target function shouldn't be call, this case should be skipped and logged 
+ if there is an explicit exception creation, skip only them, the rest should cause a fuzzing test error
+ situations that should not occur during the execution of the function should end with t.Error
+ don't use not exported fields in validation checks
	`

	patternTargetFuzz string = `
You write fuzzing test using golang

Purpose: extend target func result check in fuzzing target 

Target func:

%s

Current fuzzing target test:

%s

Current import aliases:

%s

Return new fuzzing target, that process all available func call results 

Requirements:
+ only code
+ no explanation
+ process all edge cases
+ if arguments is invalid, target function shouldn't be call, this case should be skipped and logged 
+ if there is an explicit exception creation, skip only them, the rest should cause a fuzzing test error.
+ situations that should not occur during the execution of the function should end with t.Error
+ don't use not exported fields in validation checks
	`

	patternTarget string = `
You write fuzzing test using golang

Purpose: generate fuzz target func for function:

%s

Return new fuzzing target, that process all available func call results 

Requirements:
+ only code
+ no explanation
+ process all edge cases
+ if arguments is invalid, target function shouldn't be call, this case should be skipped and logged 
+ if there is an explicit exception creation, skip only them, the rest should cause a fuzzing test error.
+ situations that should not occur during the execution of the function should end with t.Error
+ don't use not exported fields in validation checks
	`
)

func CreatePrompt(
	fset *token.FileSet,
	target *mod.Func,
	fuzzTarget string,
	qualifier *mod.ImportQualifier,
) string {
	if len(fuzzTarget) == 0 {
		return fmt.Sprintf(patternTarget, getSourceCode(target.AstFuncDecl, fset))
	}

	if len(target.Uses) > 0 {
		return fmt.Sprintf(
			patternTargetUseFuzz,
			getSourceCode(target.AstFuncDecl, fset),
			getSourceCode(target.Uses[0], fset),
			fuzzTarget,
			getImports(qualifier),
		)
	}

	return fmt.Sprintf(patternTargetFuzz, getSourceCode(target.AstFuncDecl, fset), fuzzTarget, getImports(qualifier))
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
