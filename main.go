package main

import (
	"fmt"
	"reflect"
)

func print(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

type test_t struct {
	I           *int32
	F           float32
	Str         string
	Slice_i64   []int64
	Struct      *test2_t
	Struct2     test2_t
	Slice_test3 []test3_t
	Arr_i32     [5]int32
	Arr_test3   [7]*test3_t

	notUse uint
}

type test2_t struct {
	X  int8
	Y  float64
	T3 *test3_t
}

type test3_t struct {
	Hello string
	T2    *test2_t
}

// =================================================================================

func (pOwn *test_t) clone() *test_t {
	r := new(test_t)
	r.I = new(int32)
	*r.I = *pOwn.I
	r.F = pOwn.F
	r.Str = pOwn.Str
	r.Slice_i64 = make([]int64, len(pOwn.Slice_i64))
	copy(r.Slice_i64, pOwn.Slice_i64)
	r.Struct = pOwn.Struct.clone()
	r.Struct2 = *pOwn.Struct2.clone()
	r.Slice_test3 = make([]test3_t, len(pOwn.Slice_test3))
	for i := 0; i < len(pOwn.Slice_test3); i++ {
		r.Slice_test3[i] = *pOwn.Slice_test3[i].clone()
	}
	r.Arr_i32 = pOwn.Arr_i32
	for i := 0; i < len(pOwn.Arr_test3); i++ {
		r.Arr_test3[i] = pOwn.Arr_test3[i].clone()
	}
	return r
}

func (pOwn *test2_t) clone() *test2_t {
	r := new(test2_t)
	r.X = pOwn.X
	r.Y = pOwn.Y
	r.T3 = pOwn.T3.clone()
	return r
}

func (pOwn *test3_t) clone() *test3_t {
	r := new(test3_t)
	r.Hello = pOwn.Hello
	r.T2 = pOwn.T2.clone()
	return r
}

// =================================================================================

func main() {
	gen(reflect.TypeOf(test_t{}), "r", "pOwn")
}

func gen(t reflect.Type, varPrefix string, valPrefix string) {
	var structTypes []reflect.Type
	scanStructs(t, &structTypes)

	for _, t := range structTypes {
		genStruct(t, varPrefix, valPrefix)
	}
}

func scanStructs(t reflect.Type, result *[]reflect.Type) {
	switch t.Kind() {
	case reflect.Struct:
		for _, v := range *result {
			if t == v {
				return
			}
		}

		*result = append(*result, t)
		for i := 0; i < t.NumField(); i++ {
			fld := t.Field(i)
			if !fld.IsExported() || fld.Anonymous {
				continue
			}

			scanStructs(fld.Type, result)
		}

	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		scanStructs(t.Elem(), result)
	}
}

func genStruct(t reflect.Type, varPrefix string, valPrefix string) {
	print("func (pOwn *%s) clone() *%s {", t.Name(), t.Name())
	print("%s := new(%s)", varPrefix, t.Name())
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() || fld.Anonymous {
			continue
		}

		genField(fld.Type, varPrefix, valPrefix, fld.Name)
	}
	print("return r")
	print("}")
}

func genField(t reflect.Type, varPrefix string, valPrefix string, fldName string) {
	switch t.Kind() {
	case reflect.Struct:
		print("%s.%s = *%s.%s.clone()", varPrefix, fldName, valPrefix, fldName)

	case reflect.Ptr:
		if isPrimitive(t.Elem()) {
			print("%s.%s = new(%s)", varPrefix, fldName, t.Elem().Name())
			print("*%s.%s = *%s.%s", varPrefix, fldName, valPrefix, fldName)
		} else {
			print("%s.%s = %s.%s.clone()", varPrefix, fldName, valPrefix, fldName)
		}

	case reflect.Slice, reflect.Array:
		// []int, []*int, []T , []*T
		/*
			r.fld = make([]int, len(self.fld))

			r.fld = self.fld

			copy(r.fld, self.fld)

			for i := 0; i < len(self.fld); i++ {
				r.fld[i] = self.fld[i]

				r.fld[i] = new(int)
				*r.fld[i] = *self.fld[i]

				r.fld[i] = *self.fld[i].clone()

				r.fld[i] = self.fld[i].clone()
			}
		*/
		if t.Kind() == reflect.Slice {
			print("%s.%s = make([]%s, len(%s.%s))", varPrefix, fldName, t.Elem().Name(), valPrefix, fldName)
		}

		if isPrimitive(t.Elem()) {
			if t.Kind() == reflect.Slice {
				print("copy(%s.%s, %s.%s)", varPrefix, fldName, valPrefix, fldName)
			} else {
				print("%s.%s = %s.%s", varPrefix, fldName, valPrefix, fldName)
			}
		} else {
			print("for i := 0; i < len(%s.%s); i++ {", valPrefix, fldName)
			genField(t.Elem(), varPrefix, valPrefix, fmt.Sprintf("%s[i]", fldName))
			print("}")
		}

	case reflect.Map:
	default:
		if !isPrimitive(t) {
			panic("unsupported type:" + t.String())
		}

		print("%s.%s = %s.%s", varPrefix, fldName, valPrefix, fldName)
	}
}

func isPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.String:
	default:
		return false
	}

	return true
}