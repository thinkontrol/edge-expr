package edgeexpr

import (
	"testing"
)

func TestCacheBitOperationsSupportNumericValues(t *testing.T) {
	tests := []struct {
		name  string
		value any
		bit   int
		set   bool
		byteN int
		byteI int
	}{
		{name: "uint32", value: uint32(1) << 31, bit: 31, set: true, byteN: 3, byteI: 7},
		{name: "uint64", value: uint64(1) << 63, bit: 63, set: true, byteN: 7, byteI: 7},
		{name: "bytes", value: []byte{0x05, 0x80}, bit: 15, set: true, byteN: 1, byteI: 7},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewCache(nil)
			cache.AddPoint(test.value, nil)

			got, err := cache.Bit(test.bit)
			if err != nil || got != test.set {
				t.Fatalf("Bit(%d) = %v, %v; want %v, nil", test.bit, got, err, test.set)
			}

			got, err = cache.ByteBit(test.byteN, test.byteI)
			if err != nil || got != test.set {
				t.Fatalf("ByteBit(%d, %d) = %v, %v; want %v, nil", test.byteN, test.byteI, got, err, test.set)
			}

			if got, err := cache.BitAnd(1); err != nil || got != uint(numericValue(test.value)&1) {
				t.Fatalf("BitAnd(1) = %v, %v; want %v, nil", got, err, uint(numericValue(test.value)&1))
			}
			if got, err := cache.BitOr(1); err != nil || got != uint(numericValue(test.value)|1) {
				t.Fatalf("BitOr(1) = %v, %v; want %v, nil", got, err, uint(numericValue(test.value)|1))
			}
			if got, err := cache.BitXor(1); err != nil || got != uint(numericValue(test.value)^1) {
				t.Fatalf("BitXor(1) = %v, %v; want %v, nil", got, err, uint(numericValue(test.value)^1))
			}
			if got, err := cache.BitClear(1); err != nil || got != uint(numericValue(test.value)&^1) {
				t.Fatalf("BitClear(1) = %v, %v; want %v, nil", got, err, uint(numericValue(test.value)&^1))
			}
			if test.name == "uint32" {
				if got, err := cache.BitNot(); err != nil || got != uint(^uint32(numericValue(test.value))) {
					t.Fatalf("BitNot() = %v, %v; want %v, nil", got, err, uint(^uint32(numericValue(test.value))))
				}
			}
		})
	}
}

func TestCacheBitRejectsOutOfRangeNumericBits(t *testing.T) {
	cache := NewCache(nil)
	cache.AddPoint(uint32(1), nil)

	if _, err := cache.Bit(32); err == nil {
		t.Fatal("Bit(32) should reject a uint32 value")
	}
	if _, err := cache.ByteBit(4, 0); err == nil {
		t.Fatal("ByteBit(4, 0) should reject a uint32 value")
	}
}

func numericValue(value any) uint64 {
	converted, err := bitValue(value)
	if err != nil {
		panic(err)
	}
	return converted
}
