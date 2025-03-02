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

		Example of ussing func:

		%s

		Current fuzzing target test:
	
		%s

		Return new fuzzing target, that process all available func call results 

		Requirements:
		+ only code
		+ no explanation
		+ process all edge cases
		+ if arguments is invalid, target function shouldn't be call, this case should be skipped and logged 
		+ handle only possible exceptions, the rest should be thrown above
		+ situations that should not occur during the execution of the function should end with t.Error
	`

	patternTargetFuzz string = `
		You write fuzzing test using golang

		Purpose: extend target func result check in fuzzing target 

		Target func:
		
		%s

		Current fuzzing target test:

		%s

		Return new fuzzing target, that process all available func call results 

		Requirements:
		+ only code
		+ no explanation
		+ process all edge cases
		+ if arguments is invalid, target function shouldn't be call, this case should be skipped and logged 
		+ handle only possible exceptions, the rest should be thrown above
		+ situations that should not occur during the execution of the function should end with t.Error
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
		+ handle only possible exceptions, the rest should be thrown above
		+ situations that should not occur during the execution of the function should end with t.Error
		+ as text, without code formating
	`
)

func CreatePrompt(
	fset *token.FileSet,
	target *mod.Func,
	fuzzTarget string,
) string {
	if len(fuzzTarget) == 0 {
		return fmt.Sprintf(patternTarget, GetSourceCode(target.AstFuncDecl, fset))
	}

	return fmt.Sprintf(patternTargetFuzz, GetSourceCode(target.AstFuncDecl, fset), fuzzTarget)
}

func GetSourceCode(target *ast.FuncDecl, fset *token.FileSet) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, target); err != nil {
		log.Fatalf("Failed to print function declaration: %v", err)
	}
	return buf.String()
}
