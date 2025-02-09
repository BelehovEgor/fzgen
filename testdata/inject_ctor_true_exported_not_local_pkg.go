package examplefuzz

import (
	"bufio"
	"reflect"
	"testing"

	fuzzwrapexamples "github.com/BelehovEgor/fzgen/examples/inputs/test-constructor-injection"
	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_A_PtrMethodWithArg(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *fuzzwrapexamples.A
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &i)
		if err != nil || r == nil {
			return
		}

		r.PtrMethodWithArg(i)
	})
}

func Fuzz_B_PtrMethodWithArg(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *fuzzwrapexamples.B
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &i)
		if err != nil || r == nil {
			return
		}

		r.PtrMethodWithArg(i)
	})
}

func Fuzz_MyNullUUID_UnmarshalBinary(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var nu *fuzzwrapexamples.MyNullUUID
		var data_0 []byte
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&nu, &data_0)
		if err != nil || nu == nil {
			return
		}

		nu.UnmarshalBinary(data_0)
	})
}

func Fuzz_MyRegexp_Expand(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var re *fuzzwrapexamples.MyRegexp
		var dst []byte
		var template []byte
		var src []byte
		var match []int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&re, &dst, &template, &src, &match)
		if err != nil || re == nil {
			return
		}

		re.Expand(dst, template, src, match)
	})
}

func Fuzz_Package_SetName(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pkg *fuzzwrapexamples.Package
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&pkg, &name)
		if err != nil || pkg == nil {
			return
		}

		pkg.SetName(name)
	})
}

func Fuzz_A_ValMethodWithArg(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r fuzzwrapexamples.A
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &i)
		if err != nil {
			return
		}

		r.ValMethodWithArg(i)
	})
}

func Fuzz_B_ValMethodWithArg(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r fuzzwrapexamples.B
		var i int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &i)
		if err != nil {
			return
		}

		r.ValMethodWithArg(i)
	})
}

func Fuzz_NewAPtr(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c)
		if err != nil {
			return
		}

		fuzzwrapexamples.NewAPtr(c)
	})
}

func Fuzz_NewBVal(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c)
		if err != nil {
			return
		}

		fuzzwrapexamples.NewBVal(c)
	})
}

func Fuzz_NewMyRegexp(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var a int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&a)
		if err != nil {
			return
		}

		fuzzwrapexamples.NewMyRegexp(a)
	})
}

func Fuzz_NewPackage(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var path string
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&path, &name)
		if err != nil {
			return
		}

		fuzzwrapexamples.NewPackage(path, name)
	})
}

func Fuzz_NewZ(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var z *bufio.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&z)
		if err != nil || z == nil {
			return
		}

		fuzzwrapexamples.NewZ(z)
	})
}

func fabric_interface_empty_Z(impl fuzzwrapexamples.Z) interface{} {
	return impl
}

func fabric_interface_empty_Package(impl fuzzwrapexamples.Package) interface{} {
	return impl
}

func fabric_interface_empty_MyRegexp(impl fuzzwrapexamples.MyRegexp) interface{} {
	return impl
}

func fabric_interface_empty_MyNullUUID(impl fuzzwrapexamples.MyNullUUID) interface{} {
	return impl
}

func fabric_interface_empty_B(impl fuzzwrapexamples.B) interface{} {
	return impl
}

func fabric_interface_empty_A(impl fuzzwrapexamples.A) interface{} {
	return impl
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_Z))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_Package))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_MyRegexp))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_MyNullUUID))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_B))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_A))
	FabricFuncsForCustomTypes["fuzzwrapexamples.Z"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.Z"], reflect.ValueOf(fuzzwrapexamples.NewZ))
	FabricFuncsForCustomTypes["fuzzwrapexamples.Package"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.Package"], reflect.ValueOf(fuzzwrapexamples.NewPackage))
	FabricFuncsForCustomTypes["fuzzwrapexamples.MyRegexp"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.MyRegexp"], reflect.ValueOf(fuzzwrapexamples.NewMyRegexp))
	FabricFuncsForCustomTypes["fuzzwrapexamples.B"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.B"], reflect.ValueOf(fuzzwrapexamples.NewBVal))
	FabricFuncsForCustomTypes["fuzzwrapexamples.A"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.A"], reflect.ValueOf(fuzzwrapexamples.NewAPtr))
	m.Run()
}
