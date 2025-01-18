package examplefuzz

import (
	"io"
	"reflect"
	"testing"

	fuzzwrapexamples "github.com/thepudds/fzgen/examples/inputs/test-exported"
	"github.com/thepudds/fzgen/fuzzer"
)

func Fuzz_TypeExported_PointerExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *fuzzwrapexamples.TypeExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil || t_0 == nil {
			return
		}

		t_0.PointerExportedMethod(i)
	})
}

func Fuzz_TypeExported_pointerRcvNotExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *fuzzwrapexamples.TypeExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil || t_0 == nil {
			return
		}

		t_0.pointerRcvNotExportedMethod(i)
	})
}

func Fuzz_typeNotExported_PointerExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *fuzzwrapexamples.typeNotExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil || t_0 == nil {
			return
		}

		t_0.PointerExportedMethod(i)
	})
}

func Fuzz_typeNotExported_pointerRcvNotExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *fuzzwrapexamples.typeNotExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil || t_0 == nil {
			return
		}

		t_0.pointerRcvNotExportedMethod(i)
	})
}

func Fuzz_TypeExported_NonPointerExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 fuzzwrapexamples.TypeExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil {
			return
		}

		t_0.NonPointerExportedMethod(i)
	})
}

func Fuzz_TypeExported_nonPointerRcvNotExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 fuzzwrapexamples.TypeExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil {
			return
		}

		t_0.nonPointerRcvNotExportedMethod(i)
	})
}

func Fuzz_typeNotExported_NonPointerExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 fuzzwrapexamples.typeNotExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil {
			return
		}

		t_0.NonPointerExportedMethod(i)
	})
}

func Fuzz_typeNotExported_nonPointerRcvNotExportedMethod(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 fuzzwrapexamples.typeNotExported
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &i)
		if err != nil {
			return
		}

		t_0.nonPointerRcvNotExportedMethod(i)
	})
}

func Fuzz_FuncExported(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&i)
		if err != nil {
			return
		}

		fuzzwrapexamples.FuncExported(i)
	})
}

func Fuzz_FuncExportedUsesSupportedInterface(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var w io.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&w)
		if err != nil {
			return
		}

		fuzzwrapexamples.FuncExportedUsesSupportedInterface(w)
	})
}

// skipping Fuzz_FuncExportedUsesUnsupportedInterface because parameters include unsupported type: github.com/thepudds/fzgen/examples/inputs/test-exported.ExportedInterface

func Fuzz_funcNotExported(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&i)
		if err != nil {
			return
		}

		fuzzwrapexamples.funcNotExported(i)
	})
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	m.Run()
}
