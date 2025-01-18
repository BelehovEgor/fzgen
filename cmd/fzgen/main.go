// Command fzgen automatically generates fuzz functions.
//
// See the project README for additional information:
//
//	https://github.com/BelehovEgor/fzgen
package main

import (
	"os"

	"github.com/BelehovEgor/fzgen/gen"
)

func main() {
	os.Exit(gen.FzgenMain())
}
