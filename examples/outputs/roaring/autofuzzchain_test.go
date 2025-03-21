package roaringfuzz

// Edit if desired. Code generated by "fzgen -chain -ctor=^NewBitmap$ github.com/RoaringBitmap/roaring".

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/BelehovEgor/fzgen/fuzzer"
	"github.com/RoaringBitmap/roaring"
)

func Fuzz_NewBitmap_Chain(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		fz := fuzzer.NewFuzzer(data)

		target := roaring.NewBitmap()

		steps := []fuzzer.Step{
			{
				Name: "Fuzz_Bitmap_Add",
				Func: func(x uint32) {
					target.Add(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_AddInt",
				Func: func(x int) {
					target.AddInt(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_AddMany",
				Func: func(dat []uint32) {
					target.AddMany(dat)
				},
			},
			{
				Name: "Fuzz_Bitmap_AddRange",
				Func: func(rangeStart uint64, rangeEnd uint64) {
					target.AddRange(rangeStart, rangeEnd)
				},
			},
			{
				Name: "Fuzz_Bitmap_And",
				Func: func(x2 *roaring.Bitmap) {
					target.And(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_AndAny",
				Func: func(bitmaps []*roaring.Bitmap) {
					target.AndAny(bitmaps...)
				},
			},
			{
				Name: "Fuzz_Bitmap_AndCardinality",
				Func: func(x2 *roaring.Bitmap) uint64 {
					return target.AndCardinality(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_AndNot",
				Func: func(x2 *roaring.Bitmap) {
					target.AndNot(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_CheckedAdd",
				Func: func(x uint32) bool {
					return target.CheckedAdd(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_CheckedRemove",
				Func: func(x uint32) bool {
					return target.CheckedRemove(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_Clear",
				Func: func() {
					target.Clear()
				},
			},
			{
				Name: "Fuzz_Bitmap_Clone",
				Func: func() *roaring.Bitmap {
					return target.Clone()
				},
			},
			{
				Name: "Fuzz_Bitmap_CloneCopyOnWriteContainers",
				Func: func() {
					target.CloneCopyOnWriteContainers()
				},
			},
			{
				Name: "Fuzz_Bitmap_Contains",
				Func: func(x uint32) bool {
					return target.Contains(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_ContainsInt",
				Func: func(x int) bool {
					return target.ContainsInt(x)
				},
			},
			// skipping Fuzz_Bitmap_Equals because parameters include unsupported interface: interface{}

			{
				Name: "Fuzz_Bitmap_Flip",
				Func: func(rangeStart uint64, rangeEnd uint64) {
					target.Flip(rangeStart, rangeEnd)
				},
			},
			{
				Name: "Fuzz_Bitmap_FlipInt",
				Func: func(rangeStart int, rangeEnd int) {
					target.FlipInt(rangeStart, rangeEnd)
				},
			},
			{
				Name: "Fuzz_Bitmap_Freeze",
				Func: func() ([]byte, error) {
					return target.Freeze()
				},
			},
			{
				Name: "Fuzz_Bitmap_FreezeTo",
				Func: func(buf []byte) (int, error) {
					return target.FreezeTo(buf)
				},
			},
			{
				Name: "Fuzz_Bitmap_FromBase64",
				Func: func(str string) (int64, error) {
					return target.FromBase64(str)
				},
			},
			{
				Name: "Fuzz_Bitmap_FromBuffer",
				Func: func(buf []byte) (int64, error) {
					return target.FromBuffer(buf)
				},
			},
			{
				Name: "Fuzz_Bitmap_FrozenView",
				Func: func(buf []byte) {
					target.FrozenView(buf)
				},
			},
			{
				Name: "Fuzz_Bitmap_GetCardinality",
				Func: func() uint64 {
					return target.GetCardinality()
				},
			},
			{
				Name: "Fuzz_Bitmap_GetCopyOnWrite",
				Func: func() bool {
					return target.GetCopyOnWrite()
				},
			},
			{
				Name: "Fuzz_Bitmap_GetFrozenSizeInBytes",
				Func: func() uint64 {
					return target.GetFrozenSizeInBytes()
				},
			},
			{
				Name: "Fuzz_Bitmap_GetSerializedSizeInBytes",
				Func: func() uint64 {
					return target.GetSerializedSizeInBytes()
				},
			},
			{
				Name: "Fuzz_Bitmap_GetSizeInBytes",
				Func: func() uint64 {
					return target.GetSizeInBytes()
				},
			},
			{
				Name: "Fuzz_Bitmap_HasRunCompression",
				Func: func() bool {
					return target.HasRunCompression()
				},
			},
			{
				Name: "Fuzz_Bitmap_Intersects",
				Func: func(x2 *roaring.Bitmap) bool {
					return target.Intersects(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_IsEmpty",
				Func: func() bool {
					return target.IsEmpty()
				},
			},
			// skipping Fuzz_Bitmap_Iterate because parameters include unsupported func or chan: func(x uint32) bool

			{
				Name: "Fuzz_Bitmap_Iterator",
				Func: func() roaring.IntPeekable {
					return target.Iterator()
				},
			},
			{
				Name: "Fuzz_Bitmap_ManyIterator",
				Func: func() roaring.ManyIntIterable {
					return target.ManyIterator()
				},
			},
			{
				Name: "Fuzz_Bitmap_MarshalBinary",
				Func: func() ([]byte, error) {
					return target.MarshalBinary()
				},
			},
			{
				Name: "Fuzz_Bitmap_Maximum",
				Func: func() uint32 {
					return target.Maximum()
				},
			},
			{
				Name: "Fuzz_Bitmap_Minimum",
				Func: func() uint32 {
					return target.Minimum()
				},
			},
			{
				Name: "Fuzz_Bitmap_Or",
				Func: func(x2 *roaring.Bitmap) {
					target.Or(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_OrCardinality",
				Func: func(x2 *roaring.Bitmap) uint64 {
					return target.OrCardinality(x2)
				},
			},
			{
				Name: "Fuzz_Bitmap_Rank",
				Func: func(x uint32) uint64 {
					return target.Rank(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_ReadFrom",
				Func: func(reader io.Reader, cookieHeader []byte) (int64, error) {
					return target.ReadFrom(reader, cookieHeader...)
				},
			},
			{
				Name: "Fuzz_Bitmap_Remove",
				Func: func(x uint32) {
					target.Remove(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_RemoveRange",
				Func: func(rangeStart uint64, rangeEnd uint64) {
					target.RemoveRange(rangeStart, rangeEnd)
				},
			},
			{
				Name: "Fuzz_Bitmap_ReverseIterator",
				Func: func() roaring.IntIterable {
					return target.ReverseIterator()
				},
			},
			{
				Name: "Fuzz_Bitmap_RunOptimize",
				Func: func() {
					target.RunOptimize()
				},
			},
			{
				Name: "Fuzz_Bitmap_Select",
				Func: func(x uint32) (uint32, error) {
					return target.Select(x)
				},
			},
			{
				Name: "Fuzz_Bitmap_SetCopyOnWrite",
				Func: func(val bool) {
					target.SetCopyOnWrite(val)
				},
			},
			{
				Name: "Fuzz_Bitmap_Stats",
				Func: func() roaring.Statistics {
					return target.Stats()
				},
			},
			{
				Name: "Fuzz_Bitmap_String",
				Func: func() string {
					return target.String()
				},
			},
			{
				Name: "Fuzz_Bitmap_ToArray",
				Func: func() []uint32 {
					return target.ToArray()
				},
			},
			{
				Name: "Fuzz_Bitmap_ToBase64",
				Func: func() (string, error) {
					return target.ToBase64()
				},
			},
			{
				Name: "Fuzz_Bitmap_ToBytes",
				Func: func() ([]byte, error) {
					return target.ToBytes()
				},
			},
			{
				Name: "Fuzz_Bitmap_UnmarshalBinary",
				Func: func(d1 []byte) {
					target.UnmarshalBinary(d1)
				},
			},
			{
				Name: "Fuzz_Bitmap_WriteTo",
				Func: func(stream io.Writer) (int64, error) {
					return target.WriteTo(stream)
				},
			},
			{
				Name: "Fuzz_Bitmap_Xor",
				Func: func(x2 *roaring.Bitmap) {
					target.Xor(x2)
				},
			},
		}

		// Execute a specific chain of steps, with the count, sequence and arguments controlled by fz.Chain
		fz.Chain(steps)

		// Validate with some roundtrip checks. These can be edited or deleted if not appropriate for your target.
		// Check MarshalBinary.
		result2, err := target.MarshalBinary()
		if err != nil {
			// Some targets should never return an error here for an object created by a constructor.
			// If that is the case for your target, you can change this to a panic(err) or t.Fatal.
			return
		}

		// Check UnmarshalBinary.
		var tmp2 *roaring.Bitmap
		err = tmp2.UnmarshalBinary(result2)
		if err != nil {
			panic(fmt.Sprintf("UnmarshalBinary failed after successful MarshalBinary. original: %v %#v marshalled: %q error: %v", target, target, result2, err))
		}
		if !reflect.DeepEqual(target, tmp2) {
			panic(fmt.Sprintf("MarshalBinary/UnmarshalBinary roundtrip equality failed. original: %v %#v marshalled: %q unmarshalled: %v %#v",
				target, target, result2, tmp2, tmp2))
		}
	})
}
