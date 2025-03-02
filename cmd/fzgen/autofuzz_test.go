package main

// Edit if desired. Code generated by "fzgen -llm github.com/BelehovEgor/fzgen/examples/inputs/1-supported-func".

import (
	"reflect"
	"runtime"
	"testing"

	supportedfunc "github.com/BelehovEgor/fzgen/examples/inputs/1-supported-func"
	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_N2_Log(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var emit supportedfunc.EmitFunc
		var lvl int
		var log string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&emit, &lvl, &log)
		if err != nil {
			t.Logf("Skipping invalid input: %!v(MISSING)", err)
			return
		}

		if emit == nil {
			t.Error("Emit function cannot be nil")
			return
		}

		if lvl < 0 || lvl > 5 {
			t.Logf("Skipping invalid log level: %!d(MISSING)", lvl)
			return
		}

		if log == "" {
			t.Logf("Skipping empty log message")
			return
		}

		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					t.Errorf("Recovered from panic: %!v(MISSING)", r)
				default:
					t.Errorf("Recovered from unexpected panic: %!v(MISSING)", r)
				}
			}
		}()

		supportedfunc.Log(emit, lvl, log)
	})
}

func Fuzz_N3_LogV2(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var emit func(lvl int, format string, args ...interface{})
		var lvl int
		var log string

		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&emit, &lvl, &log)
		if err != nil {
			t.Logf("Skipping invalid input: %!v(MISSING)", err)
			return
		}

		if emit == nil {
			t.Error("emit function cannot be nil")
			return
		}

		if lvl < 0 {
			t.Error("lvl cannot be negative")
			return
		}

		if log == "" {
			t.Error("log string cannot be empty")
			return
		}

		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case error:
					t.Error(r)
				default:
					t.Errorf("Recovered from unexpected panic: %!v(MISSING)", r)
				}
			}
		}()

		supportedfunc.LogV2(emit, lvl, log)
	})
}

func Fuzz_N4_MyEmitter(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var lvl int
		var format string
		var args []interface{}
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes, t, fuzzer.Constructors)
		err := fz.Fill2(&lvl, &format, &args)
		if err != nil {
			t.Logf("Skipping invalid input: %!v(MISSING)", err)
			return
		}

		if format == "" {
			t.Logf("Skipping empty format string")
			return
		}

		if len(args) == 0 {
			t.Logf("Skipping no arguments provided")
			return
		}

		for _, arg := range args {
			if arg == nil {
				t.Logf("Skipping nil argument")
				return
			}
		}

		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(runtime.Error); ok {
					panic(r)
				}
				t.Errorf("Recovered from panic: %!v(MISSING)", r)
			}
		}()

		supportedfunc.MyEmitter(lvl, format, args...)
	})
}
func fabric_interface_empty_string(impl string) interface{} {
	return impl
}

func fabric_func_2() func(int, string, ...interface{}) {
	return supportedfunc.MyEmitter
}

func fabric_func_1() supportedfunc.EmitFunc {
	return supportedfunc.MyEmitter
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	FabricFuncsForCustomTypes["supportedfunc.EmitFunc"] = append(FabricFuncsForCustomTypes["supportedfunc.EmitFunc"], reflect.ValueOf(fabric_func_1))
	FabricFuncsForCustomTypes["func(int,string,...interface{})"] = append(FabricFuncsForCustomTypes["func(int,string,...interface{})"], reflect.ValueOf(fabric_func_2))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_string))
	m.Run()
}
