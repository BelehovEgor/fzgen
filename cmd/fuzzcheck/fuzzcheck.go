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
	"regexp"
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
	interCount := 0
	totalCount := 0
	for _, target := range fuzzTargets {
		success, inter, total := runFuzzTarget(filename, target, 60*time.Second)
		interCount += inter
		totalCount += total

		if success {
			successCount++
		}
	}

	// Step 5: Report results
	fmt.Printf("File is valid: true\n")
	fmt.Printf("Successful fuzz tests: %d/%d; new interesting/total %d/%d\n", successCount, len(fuzzTargets), interCount, totalCount)
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
func runFuzzTarget(filename, target string, duration time.Duration) (bool, int, int) {
	dir := filepath.Dir(filename)
	cmd := exec.Command("go", "test", "-fuzz", fmt.Sprintf("^%s$", target), "-fuzztime", duration.String())
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	fmt.Printf("Running fuzz target %s for %s...\n", target, duration)
	err := cmd.Run()

	interesting, total := getLastInterestingAndTotal(out.String())

	if err != nil {
		fmt.Printf("Fuzz target %s failed: %s\n", target, err)
		fmt.Println(out.String())
		cleanupError()
		return false, interesting, total
	}

	fmt.Printf("Fuzz target %s succeeded.\n", target)

	return true, interesting, total
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

func getLastInterestingAndTotal(log string) (int, int) {
	// Define a regular expression to match the "new interesting" and "total" values
	re := regexp.MustCompile(`new interesting: (\d+) \(total: (\d+)\)`)

	// Split the log into lines
	lines := strings.Split(log, "\n")

	// Iterate over the lines in reverse order to find the last occurrence
	for i := len(lines) - 1; i >= 0; i-- {
		matches := re.FindStringSubmatch(lines[i])
		if len(matches) == 3 {
			// Extract the "new interesting" and "total" values
			newInteresting := matches[1]
			total := matches[2]

			// Convert the extracted strings to integers
			var newInterestingInt, totalInt int
			_, err := fmt.Sscanf(newInteresting, "%d", &newInterestingInt)
			if err != nil {
				return 0, 0
			}
			_, err = fmt.Sscanf(total, "%d", &totalInt)
			if err != nil {
				return 0, 0
			}

			return newInterestingInt, totalInt
		}
	}

	// If no match is found, return an error
	return 0, 0
}
