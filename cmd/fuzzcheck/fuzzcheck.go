package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run script.go <gofile>")
	}
	filename := os.Args[1]

	// Step 1: Read the Go file
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Step 2: Check if the file is valid Go code
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, filename, content, parser.AllErrors)
	if err != nil {
		log.Fatalf("File is not valid Go code: %s", err)
	}
	fmt.Println("File is valid Go code.")

	// Step 3: Get all fuzz targets
	fuzzTargets := getFuzzTargets(content)
	if len(fuzzTargets) == 0 {
		log.Fatal("No fuzz targets found.")
	}
	fmt.Printf("Found %d fuzz targets: %v\n", len(fuzzTargets), fuzzTargets)

	// Step 4: Run each fuzz target for 10 seconds
	successCount := 0
	for _, target := range fuzzTargets {
		if runFuzzTarget(filename, target, 10*time.Second) {
			successCount++
		}
	}

	// Step 5: Report results
	fmt.Printf("File is valid: true\n")
	fmt.Printf("Successful fuzz tests: %d/%d\n", successCount, len(fuzzTargets))
}

// getFuzzTargets extracts all fuzz target function names from the Go file content.
func getFuzzTargets(content []byte) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}

	var targets []string
	ast.Inspect(file, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if strings.HasPrefix(fn.Name.Name, "Fuzz") && len(fn.Type.Params.List) == 1 {
			targets = append(targets, fn.Name.Name)
		}
		return true
	})
	return targets
}

// runFuzzTarget runs a specific fuzz target for a given duration.
func runFuzzTarget(filename, target string, duration time.Duration) bool {
	dir := filepath.Dir(filename)
	cmd := exec.Command("go", "test", "-fuzz", fmt.Sprintf("^%s$", target), "-fuzztime", duration.String())
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	fmt.Printf("Running fuzz target %s for %s...\n", target, duration)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Fuzz target %s failed: %s\n", target, err)
		fmt.Println(out.String())
		cleanupError()
		return false
	}
	fmt.Printf("Fuzz target %s succeeded.\n", target)
	return true
}

func cleanupError() {
	// Path to the testdata directory
	testdataDir := "testdata"

	// Remove the testdata directory and its contents
	err := os.RemoveAll(testdataDir)
	if err != nil {
		fmt.Printf("Failed to remove testdata directory: %v\n", err)
	} else {
		fmt.Println("Successfully removed testdata directory")
	}
}
