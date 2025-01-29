package examplefuzz

import (
	"context"
	"io"
	"reflect"
	"testing"
	"unsafe"

	fuzzwrapexamples "github.com/BelehovEgor/fzgen/examples/inputs/test-types"
	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_TypesNilCheck_Interface(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var n *fuzzwrapexamples.TypesNilCheck
		var x1 io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&n, &x1)
		if err != nil || n == nil {
			return
		}

		n.Interface(x1)
	})
}

func Fuzz_TypesNilCheck_Pointers(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var n *fuzzwrapexamples.TypesNilCheck
		var x1 *int
		var x2 **int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&n, &x1, &x2)
		if err != nil || n == nil || x1 == nil || x2 == nil {
			return
		}

		n.Pointers(x1, x2)
	})
}

func Fuzz_TypesNilCheck_WriteTo(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var n *fuzzwrapexamples.TypesNilCheck
		var stream io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&n, &stream)
		if err != nil || n == nil {
			return
		}

		n.WriteTo(stream)
	})
}

func Fuzz_Std_ListenPacket(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var std_ fuzzwrapexamples.Std
		var ctx context.Context
		var network string
		var address string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&std_, &ctx, &network, &address)
		if err != nil {
			return
		}

		std_.ListenPacket(ctx, network, address)
	})
}

func Fuzz_Discard(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var string_ string
		var _arr_interface_ []interface{}
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&string_, &_arr_interface_)
		if err != nil {
			return
		}

		fuzzwrapexamples.Discard(string_, _arr_interface_...)
	})
}

func Fuzz_Discard2(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var string_ string
		var _arr_int_ []int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&string_, &_arr_int_)
		if err != nil {
			return
		}

		fuzzwrapexamples.Discard2(string_, _arr_int_...)
	})
}

func Fuzz_InterfacesFullList(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 io.Writer
		var x2 io.Reader
		var x3 io.ReaderAt
		var x4 io.WriterTo
		var x5 io.Seeker
		var x6 io.ByteScanner
		var x7 io.RuneScanner
		var x8 io.ReadSeeker
		var x9 io.ByteReader
		var x10 io.RuneReader
		var x11 io.ByteWriter
		var x12 io.ReadWriter
		var x13 io.ReaderFrom
		var x14 io.StringWriter
		var x15 io.Closer
		var x16 io.ReadCloser
		var x17 context.Context
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1, &x2, &x3, &x4, &x5, &x6, &x7, &x8, &x9, &x10, &x11, &x12, &x13, &x14, &x15, &x16, &x17)
		if err != nil {
			return
		}

		fuzzwrapexamples.InterfacesFullList(x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13, x14, x15, x16, x17)
	})
}

func Fuzz_InterfacesShortList(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var ctx context.Context
		var w io.Writer
		var r io.Reader
		var sw io.StringWriter
		var rc io.ReadCloser
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&ctx, &w, &r, &sw, &rc)
		if err != nil {
			return
		}

		fuzzwrapexamples.InterfacesShortList(ctx, w, r, sw, rc)
	})
}

// skipping Fuzz_InterfacesSkip because parameters include unsupported type: net.Conn

func Fuzz_Short1(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short1(x1)
	})
}

func Fuzz_Short2(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 *int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil || x1 == nil {
			return
		}

		fuzzwrapexamples.Short2(x1)
	})
}

func Fuzz_Short3(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 **int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil || x1 == nil {
			return
		}

		fuzzwrapexamples.Short3(x1)
	})
}

func Fuzz_Short4(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 fuzzwrapexamples.MyInt
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short4(x1)
	})
}

func Fuzz_Short5(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 complex64
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short5(x1)
	})
}

func Fuzz_Short6(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 complex128
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short6(x1)
	})
}

func Fuzz_Short7(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 uintptr
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short7(x1)
	})
}

func Fuzz_Short8(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 unsafe.Pointer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1)
		if err != nil {
			return
		}

		fuzzwrapexamples.Short8(x1)
	})
}

func Fuzz_TypesShortListFill(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 int
		var x2 *int
		var x3 **int
		var x4 map[string]string
		var x5 *map[string]string
		var x6 fuzzwrapexamples.MyInt
		var x7 [4]int
		var x8 fuzzwrapexamples.MyStruct
		var x9 io.ByteReader
		var x10 io.RuneReader
		var x11 io.ByteWriter
		var x12 io.ReadWriter
		var x13 io.ReaderFrom
		var x14 io.StringWriter
		var x15 io.Closer
		var x16 io.ReadCloser
		var x17 context.Context
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1, &x2, &x3, &x4, &x5, &x6, &x7, &x8, &x9, &x10, &x11, &x12, &x13, &x14, &x15, &x16, &x17)
		if err != nil || x2 == nil || x3 == nil || x5 == nil {
			return
		}

		fuzzwrapexamples.TypesShortListFill(x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13, x14, x15, x16, x17)
	})
}

func Fuzz_TypesShortListNoFill(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x1 int
		var x5 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x1, &x5)
		if err != nil {
			return
		}

		fuzzwrapexamples.TypesShortListNoFill(x1, x5)
	})
}

// skipping Fuzz_TypesShortListSkip1 because parameters include unsupported type: chan bool

func Fuzz_TypesShortListSkip2(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var x func(int)
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&x)
		if err != nil {
			return
		}

		fuzzwrapexamples.TypesShortListSkip2(x)
	})
}

func fabric_interface_empty_TypesNilCheck(impl fuzzwrapexamples.TypesNilCheck) interface{} {
	return impl
}

func fabric_interface_empty_Std(impl fuzzwrapexamples.Std) interface{} {
	return impl
}

func fabric_interface_empty_MyStruct(impl fuzzwrapexamples.MyStruct) interface{} {
	return impl
}

func fabric_func_21_0() func(int) {
	return fuzzwrapexamples.Short1
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_TypesNilCheck))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_Std))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty_MyStruct))
	FabricFuncsForCustomTypes["func(int)"] = append(FabricFuncsForCustomTypes["func(int)"], reflect.ValueOf(fabric_func_21_0))
	FabricFuncsForCustomTypes["fuzzwrapexamples.TypesNilCheck"] = append(FabricFuncsForCustomTypes["fuzzwrapexamples.TypesNilCheck"], reflect.ValueOf(fuzzwrapexamples.NewTypesNilCheck))
	m.Run()
}
