package main

import (
	"reflect"
	"unsafe"
)

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if x.Type() != y.Type() {
		return false
	}

	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()
	case reflect.String:
		return x.String() == y.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return x.Int() == y.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()
	case reflect.Float32, reflect.Float64:
		return equalFloats(x.Float(), y.Float())
	case reflect.Complex64, reflect.Complex128:
		return equalComplexes(x.Complex(), y.Complex())
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return x.Pointer() == y.Pointer()
	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)
	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true
	case reflect.Struct:
		if x.NumField() != y.NumField() {
			return false
		}
		for i, n := 0, x.NumField(); i < n; i++ {
			if !equal(x.Field(i), y.Field(i), seen) {
				return false
			}
		}
		return true
	case reflect.Map:
		if x.Len() != y.Len() {
			return false
		}
		for _, k := range x.MapKeys() {
			if !equal(x.MapIndex(k), y.MapIndex(k), seen) {
				return false
			}
		}
		return true
	}
	return false
}

const FloatDiff = 1.0e-10

func equalFloats(x, y float64) bool {
	lowX := x - FloatDiff
	highX := x + FloatDiff
	return lowX <= y && highX >= y
}

func equalComplexes(x, y complex128) bool {
	return equalFloats(real(x), real(y)) && equalFloats(imag(x), imag(y))
}

func Equal(x, y any) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

func IsCycle(x any) bool {
	seen := make([]unsafe.Pointer, 0)
	return isCycle(reflect.ValueOf(x), seen)
}

func isCycle(x reflect.Value, seen []unsafe.Pointer) bool {
	if !x.IsValid() {
		return false
	}
	if !x.CanAddr() &&
		x.Kind() != reflect.Struct &&
		x.Kind() != reflect.Array {

		xptr := unsafe.Pointer(x.UnsafeAddr())
		for _, s := range seen {
			if xptr == s {
				return true
			}
		}
		seen = append(seen, xptr)
	}
	switch x.Kind() {
	case reflect.Ptr, reflect.Interface:
		return isCycle(x.Elem(), seen)
	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if isCycle(x.Field(i), seen) {
				return true
			}
		}
		return false
	case reflect.Slice, reflect.Array:
		for i := 0; i < x.Len(); i++ {
			if isCycle(x.Index(i), seen) {
				return true
			}
		}
		return false
	case reflect.Map:
		for _, k := range x.MapKeys() {
			if isCycle(x.MapIndex(k), seen) {
				return true
			}
		}
		return false
	}
	return false
}
