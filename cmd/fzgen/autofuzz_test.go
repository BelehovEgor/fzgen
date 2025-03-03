package main

// Edit if desired. Code generated by "fzgen --llm=openrouter github.com/BelehovEgor/fzgen/examples/inputs/1-supported-func".

import (
	"reflect"
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
			return
		}

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
			return
		}

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
			return
		}

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
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_string))
	FabricFuncsForCustomTypes["supportedfunc.EmitFunc"] = append(FabricFuncsForCustomTypes["supportedfunc.EmitFunc"], reflect.ValueOf(fabric_func_1))
	FabricFuncsForCustomTypes["func(int,string,...interface{})"] = append(FabricFuncsForCustomTypes["func(int,string,...interface{})"], reflect.ValueOf(fabric_func_2))
	m.Run()
}
