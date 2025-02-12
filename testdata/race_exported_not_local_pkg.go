package examplefuzz

import (
	"testing"

	raceexample "github.com/BelehovEgor/fzgen/examples/inputs/race"
	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_NewMySafeMap_Chain(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		fz := fuzzer.NewFuzzer(data)

		target := raceexample.NewMySafeMap()

		steps := []fuzzer.Step{
			{
				Name: "Fuzz_MySafeMap_Load",
				Func: func(key [16]byte) *raceexample.Request {
					return target.Load(key)
				},
			},
			{
				Name: "Fuzz_MySafeMap_Store",
				Func: func(key [16]byte, req *raceexample.Request) {
					target.Store(key, req)
				},
			},
		}

		// Execute a specific chain of steps, with the count, sequence and arguments controlled by fz.Chain
		fz.Chain(steps, fuzzer.ChainParallel)
	})
}
