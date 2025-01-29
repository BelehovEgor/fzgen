package examplefuzz

import (
	"bufio"
	"reflect"
	"testing"

	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_A_PtrMethodWithArg(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *A
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
		var r *B
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
		var nu *MyNullUUID
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
		var re *MyRegexp
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
		var pkg *Package
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
		var r A
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
		var r B
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

		NewAPtr(c)
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

		NewBVal(c)
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

		NewMyRegexp(a)
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

		NewPackage(path, name)
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

		NewZ(z)
	})
}

func fabric_interface_empty_Z(impl Z) interface{} {
	return impl
}

func fabric_interface_empty_Package(impl Package) interface{} {
	return impl
}

func fabric_interface_empty_MyRegexp(impl MyRegexp) interface{} {
	return impl
}

func fabric_interface_empty_MyNullUUID(impl MyNullUUID) interface{} {
	return impl
}

func fabric_interface_empty_B(impl B) interface{} {
	return impl
}

func fabric_interface_empty_A(impl A) interface{} {
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
	FabricFuncsForCustomTypes["Z"] = append(FabricFuncsForCustomTypes["Z"], reflect.ValueOf(NewZ))
	FabricFuncsForCustomTypes["Package"] = append(FabricFuncsForCustomTypes["Package"], reflect.ValueOf(NewPackage))
	FabricFuncsForCustomTypes["MyRegexp"] = append(FabricFuncsForCustomTypes["MyRegexp"], reflect.ValueOf(NewMyRegexp))
	FabricFuncsForCustomTypes["B"] = append(FabricFuncsForCustomTypes["B"], reflect.ValueOf(NewBVal))
	FabricFuncsForCustomTypes["A"] = append(FabricFuncsForCustomTypes["A"], reflect.ValueOf(NewAPtr))
	m.Run()
}
