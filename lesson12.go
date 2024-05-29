package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

var nextLevel int

func Display(name string, x any) {
	fmt.Printf("Display %s (%T)\n", name, x)
	nextLevel = 0
	display(name, reflect.ValueOf(x))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Map, reflect.Slice:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Struct:
		var b bytes.Buffer
		b.WriteString(v.Type().String())
		b.WriteRune('{')
		for i := 0; i < v.NumField(); i++ {
			b.WriteString(fmt.Sprintf("%s: %s", v.Type().Field(i).Name, formatAtom(v.Field(i))))
			if i != v.NumField()-1 {
				b.WriteString(", ")
			}
		}
		b.WriteRune('}')
		return b.String()
	case reflect.Array:
		var b bytes.Buffer
		b.WriteString(v.Type().String())
		b.WriteRune('[')
		for i := 0; i < v.Len(); i++ {
			b.WriteString(formatAtom(v.Index(i)))
			if i != v.Len()-1 {
				b.WriteString(", ")
			}
		}
		b.WriteRune(']')
		return b.String()
	default:
		return v.Type().String() + " value"
	}
}

func display(path string, x reflect.Value) {
	nextLevel++
	if nextLevel > 20 {
		return
	}

	switch x.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s is invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < x.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), x.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, x.Type().Field(i).Name)
			display(fieldPath, x.Field(i))
		}
	case reflect.Map:
		for _, k := range x.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(k)), x.MapIndex(k))
		}
	case reflect.Ptr:
		if x.IsNil() {
			fmt.Printf("%s is nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, x.Elem())
			display(path+".value", x.Elem())
		}
	case reflect.Interface:
		if x.IsNil() {
			fmt.Printf("%s is nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, x.Elem().Type())
			display(path+".value", x.Elem())
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(x))
	}
}

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("invalid")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, _ = fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, _ = fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		_, _ = fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			_, _ = fmt.Fprintf(buf, "(%s", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Bool:
		if v.Bool() {
			_, _ = fmt.Fprintf(buf, "t")
		} else {
			_, _ = fmt.Fprintf(buf, "nil")
		}
	case reflect.Float32, reflect.Float64:
		_, _ = fmt.Fprintf(buf, "%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		v := v.Complex()
		_, _ = fmt.Fprintf(buf, "$C(%f %f)", real(v), imag(v))
	case reflect.Interface:
		buf.WriteByte('(')
		t := v.Type()
		if t.Name() == "" {
			_, _ = fmt.Fprintf(buf, "%q ", v.Elem().Type().String())
		} else {
			_, _ = fmt.Fprintf(buf, "%s.%s", t.PkgPath(), t.Name())
		}
		if err := encode(buf, v.Elem()); err != nil {
			return err
		}
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Unmarshal(data []byte, out any) (err error) {

}
