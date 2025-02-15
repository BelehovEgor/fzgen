// Package gen automatically generates fuzz functions, and is the main entry point
// for the fzgen command.
//
// See the project README for additional information:
//
//	https://github.com/BelehovEgor/fzgen
package gen

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/tools/imports"
)

// one way to test this is against the stdlib (here, this just tests that fzgen generates and the result compiles successfully):
//   for x in $(go list std | egrep -v 'internal|runtime|unsafe|vendor|image/color/palette'); do echo $x; mkdir -p $x; cd $x ; fzgen $x && gotip test -c || echo "--- FAILED $x ---"; cd - >/dev/null; done &> out.txt
// current stats:
//   grep -r '^func Fuzz' | wc -l
//   2775
//   grep -r 'skipping' | wc -l
//   603

// Usage contains short usage information.
var Usage = `
Usage:
	fzgen [-chain] [-parallel] [-ctor=<target-constructor-regexp>] [-unexported] [-mocks] [-mocksPackagePrefix] [-package] [-mocksDepth] [-prioritizeConstructors]
	
Running fzgen without any arguments targets the package in the current directory.

fzgen outputs a set of wrapper fuzz functions for all functions matching
the -func regexp, which defaults to matching all functions. The target package
defaults to the current directory. The target package must be in the current
module or listed as dependency of the current module (e.g., via 'go get example.com/foo').

The resulting wrapper functions will all start with 'Fuzz', and are candidates 
for use with fuzzing via Go 1.18 cmd/go (e.g., 'gotip test -fuzz=.').

fzgen supports a package pattern as the last argument. If multiple packages match,
the generated files will be placed in each target package's directory. Otherwise,
when there is only a single target package, the generated file will be placed
in the current working directory.

Test functions and any function that already starts with 'Fuzz' are skipped,
as are functions that have unsupported parameters such as a channel.

`

var (
	debugForceLocal        bool
	debugForceNoTopComment bool
)

func FzgenMain() int {
	// handle flags
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, Usage)
		flag.PrintDefaults()
	}

	// Most commonly used:
	chainFlag := flag.Bool("chain", false, "loop over the methods of an object, which requires finding a suitable constructor in the same package and which is controllable via the -ctor flag.")
	parallelFlag := flag.Bool("parallel", false, "indicates an emitted chain can be run in parallel. requires -chain")
	outFileFlag := flag.String("o", "autofuzz_test.go", "output file name. defaults to autofuzz_test.go or autofuzzchain_test.go")
	constructorPatternFlag := flag.String("ctor", ".", "regexp to use if searching for constructors to automatically use.")

	// Less commonly used:
	funcPatternFlag := flag.String("func", ".", "function regex, defaults to matching all candidate functions")
	unexportedFlag := flag.Bool("unexported", false, "emit wrappers for unexported functions in addition to exported functions")

	// Mocks
	mocksEnabled := flag.Bool("mocks", false, "generates mocks for all interfaces that can be argument of functions")
	mocksPackagePrefix := flag.String("mocksPackagePrefix", "", "the package of the module from which fzgen is launche ({moduleName}/{folder struct})")
	maxDepth := flag.Int("mocksDepth", 3, "max mock depth (default: 3)")

	prioritizeConstructors := flag.Bool("prioritizeConstructors", true, "Use constructors to populate structures in priority")

	flag.Parse()

	var pkgPattern string
	switch {
	case flag.NArg() > 1:
		fmt.Println("fzgen: only one package pattern argument is allowed, and it must match only one package")
		flag.Usage()
		return 2
	case flag.NArg() == 1:
		pkgPattern = flag.Arg(0)
	default:
		pkgPattern = "."
	}

	// Search for functions in the requested packages that match the supplied func and ctor regex.
	options := flagExcludeFuzzPrefix | flagMultiMatch
	if !*unexportedFlag {
		options |= flagRequireExported
	}

	if *parallelFlag && !*chainFlag {
		fmt.Fprint(os.Stderr, "fzgen: -parallel flag requires -chain\n")
		return 2
	}

	if *mocksEnabled && *mocksPackagePrefix == "" {
		fmt.Fprint(os.Stderr, "fzgen: -mocks flag requires -mocksPackagePrefix\n")
		return 2
	}

	if *chainFlag && *outFileFlag == "autofuzz_test.go" {
		*outFileFlag = "autofuzzchain_test.go"
	}

	analyzeResult, err := findFuncsGrouped(pkgPattern, *funcPatternFlag, *constructorPatternFlag, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fzgen: %v\n", err)
		return 1
	}

	pkgs := analyzeResult.Packages

	// Check if we are looking at one package vs. multiple.
	if len(pkgs) > 1 && hasPath(*outFileFlag) {
		fmt.Fprint(os.Stderr, "fzgen:-o can only specify a file name and not a path when the package pattern matches multiple packages\n")
		return 2
	}

	wd, err := os.Getwd()
	if err != nil {
		fail(err)
	}

	var topComment string
	if !debugForceNoTopComment {
		args := os.Args[1 : len(os.Args)-flag.NArg()]
		args = append(args, []string{pkgPattern}...)
		topComment = fmt.Sprintf(topCommentTmpl, strings.Join(args, " "))
	}

	// Loop over our packages, and start our real work.
	var generatedFiles int
	for i := range pkgs {
		if len(pkgs[i].Targets) == 0 {
			continue
		}

		// Determine what output file we will create, and what package name we will use in it.
		var wrapperPkgName string
		outFile := *outFileFlag
		targetPkgName := pkgs[i].Targets[0].PkgName
		switch {
		case len(pkgs) > 1:
			// Specifying multiple packages via a pattern creates an output file in each package's directory.
			outFile = filepath.Join(pkgs[i].Targets[0].PkgDir, outFile)
			wrapperPkgName = targetPkgName
		case len(pkgs) == 1:
			// When the target is a single package overall, we default to placing the output file in the working directory,
			// with the ability to set a more specific output file path via -o. For the working directory or some other user-supplied
			// output location, there might already be a package there, so we need to look up that package name if it exists.
			wrapperPkgName = outDirPkgName(outFile)
			if wrapperPkgName == "" {
				// We did not find a package name in our destination (e.g., might not have any .go files there),
				// so make up a new package name.
				wrapperPkgName = targetPkgName + "fuzz"
			}
		default:
			panic("impossible")
		}

		// If the output file will end up in the target package dir, we set qualifyAll to false
		// so that emitted references to types from the target package will not be qualified with the package prefix.
		outDir, err := filepath.Abs(filepath.Dir(outFile))
		if err != nil {
			fail(err)
		}
		targetDir, err := filepath.Abs(pkgs[i].Targets[0].PkgDir)
		if err != nil {
			fail(err)
		}
		qualifyAll := targetDir != outDir

		if debugForceLocal {
			qualifyAll = false
		}

		wrapperOpts := wrapperOptions{
			qualifyAll:                           qualifyAll,
			topComment:                           topComment,
			requiredMocks:                        *mocksEnabled,
			mocksPackagePrefix:                   *mocksPackagePrefix,
			maxMockDepth:                         *maxDepth,
			constructorsArePriorityForStructures: *prioritizeConstructors,
		}

		// Do the actual work of emitting our wrappers.
		var gen *generated
		if !*chainFlag {
			gen, err = emitIndependentWrappers(outDir, pkgs[i], analyzeResult.TypeContext, wrapperPkgName, wrapperOpts)
		} else {
			gen, err = emitChainWrappers(outDir, pkgs[i], analyzeResult.TypeContext, wrapperPkgName, wrapperOpts)
		}

		// Handle certain common errors gracefully, including skipping & continuing if multiple target packages.
		if err != nil {
			fail(err)
		}

		generatedFiles++

		if gen.MockeryYaml != nil {
			err = os.WriteFile(fmt.Sprintf("%s/.mockery.yaml", outDir), gen.MockeryYaml, 0o644)
			if err != nil {
				fail(err)
			}

			cmd := exec.Command("mockery")
			cmd.Dir = outDir
			_, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				fail(err)
			}
		}

		//Fix up any needed imports.
		var adjusted []byte
		abs, err := filepath.Abs(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fzgen: warning: continuing after failing to find abs path:", err)
			abs = outFile
		}
		adjusted, err = imports.Process(abs, gen.Tests, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "fzgen: warning: continuing after failing to automatically adjust imports:", err)
			adjusted = gen.Tests
		}

		// Write the output.
		err = os.WriteFile(outFile, adjusted, 0o644)
		if err != nil {
			fail(err)
		}

		rel := outFile
		if hasPath(outFile) {
			rel, err = filepath.Rel(wd, outFile)
			if err != nil {
				rel = outFile
			}
		}
		fmt.Println("fzgen: created", rel)
	}

	switch generatedFiles {
	case 0:
		fmt.Printf("fzgen: no targets were found\n")
	case 1:
	default:
		fmt.Printf("fzgen: created %d files\n", generatedFiles)
	}

	return 0
}

// hasPath reports whether the string includes any path elements.
// hasPath treats both forward slash  and back slash as indicating a path,
// which is somewhat similar to similar to '-o' flag for go build looking
// for both varieties of slash.
func hasPath(s string) bool {
	if strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return true
	}
	cleaned := filepath.Clean(s)
	if cleaned == "." || cleaned == ".." {
		return true
	}
	if filepath.Dir(s) != "." {
		return true
	}
	return false
}

// outDirPkgName determines the package name for the directory that
// contains outFile, or returns the empty string if it does not find an
// existing package in that directory.
func outDirPkgName(outFile string) string {
	outDir, err := filepath.Abs(filepath.Dir(outFile))
	if err != nil {
		fail(err)
	}

	isMod, err := isInModule(outDir)
	if err != nil {
		fail(fmt.Errorf("failed when checking if directory %q is a module: %v", outDir, err))
	}
	if !isMod {
		fail(fmt.Errorf("output directory %q is not a module", outDir))
	}

	// Determine our current package name using go list.
	pkgNames, err := goList(outDir, "-e", "-f", "{{.Name}}", ".")
	if err != nil {
		fail(err)
	}
	switch len(pkgNames) {
	case 0:
		// No .go files or possibly no valid .go files in the current dir.
		return ""
	case 1:
		return pkgNames[0]
	default:
		fail(errors.New("unexpected"))
	}
	return ""
}

// fail prints an error to stderr and exits.
func fail(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		file = filepath.Base(file)
		fmt.Fprintf(os.Stderr, "fzgen: %s:%d: error: %v\n", file, line, err)
	} else {
		fmt.Fprintf(os.Stderr, "fzgen: error: %v\n", err)
	}
	os.Exit(1)
}

func init() {
	debug := strings.Split(os.Getenv("FZDEBUG"), ",")
	for _, f := range debug {
		if strings.HasPrefix(f, "forcelocal=") {
			debugVal, err := strconv.Atoi(strings.TrimPrefix(f, "forcelocal="))
			if err != nil || debugVal > 1 {
				panic("unexpected forcelocal value in FZDEBUG env var")
			}
			if debugVal == 1 {
				debugForceLocal = true
			}
		}
		if strings.HasPrefix(f, "notopcomment=") {
			debugVal, err := strconv.Atoi(strings.TrimPrefix(f, "notopcomment="))
			if err != nil || debugVal > 1 {
				panic("unexpected notopcomment value in FZDEBUG env var")
			}
			if debugVal == 1 {
				debugForceNoTopComment = true
			}
		}
	}
}

var topCommentTmpl string = `
// Edit if desired. Code generated by "fzgen %s".

`
