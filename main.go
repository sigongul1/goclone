package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
)

// =================================================================================

func main() {
	gen(reflect.TypeOf(test_t{}))
	// example_check()
}

// =================================================================================

var iNames = []string{"i", "j", "k", "z1", "z2", "z3"}
var kNames = []string{"k", "k2", "k3", "k4", "k5"}
var vNames = []string{"v", "v2", "v3", "v4", "v5"}

func print(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func gen(t reflect.Type) {
	// scan structs
	structs := make(map[reflect.Type]struct{})
	scanStructs(t, structs, nil)

	// sort
	var structTypes []reflect.Type
	for t := range structs {
		structTypes = append(structTypes, t)
	}
	sort.Slice(structTypes, func(i, j int) bool {
		return structTypes[i].Name() < structTypes[j].Name()
	})

	// gen
	for _, t := range structTypes {
		genStruct(t, "r", "pOwn", 0, 0)
	}
}

func scanStructs(t reflect.Type, result map[reflect.Type]struct{}, tracker []reflect.Type) {
	switch t.Kind() {
	case reflect.Struct:
		for _, v := range tracker {
			if t == v {
				print("ERROR: recursive references")
				os.Exit(1)
			}
		}
		tracker = append(tracker, t)
		result[t] = struct{}{}

		for i := 0; i < t.NumField(); i++ {
			fld := t.Field(i)
			if !fld.IsExported() || fld.Anonymous {
				continue
			}

			scanStructs(fld.Type, result, tracker)
		}

	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		scanStructs(t.Elem(), result, tracker)
	}
}

func genStruct(t reflect.Type, varStr string, valStr string, iDepth, kvDepth int) {
	print("func (pOwn *%s) clone() *%s {", t.Name(), t.Name())
	print("%s := new(%s)\n", varStr, t.Name())
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !fld.IsExported() || fld.Anonymous {
			continue
		}

		genField(
			fld.Type,
			fmt.Sprintf("%s.%s", varStr, fld.Name),
			fmt.Sprintf("%s.%s", valStr, fld.Name),
			iDepth, kvDepth,
		)
	}
	print("\nreturn %s", varStr)
	print("}\n")
}

func genField(t reflect.Type, varStr string, valStr string, iDepth, kvDepth int) {
	switch t.Kind() {
	case reflect.Struct:
		print("%s = *%s.clone()", varStr, valStr)

	case reflect.Ptr:
		print("if %s != nil {", valStr)
		if isPrimitive(t.Elem()) {
			print("%s = new(%s)", varStr, getTypeName(t.Elem()))
			print("*%s = *%s", varStr, valStr)
		} else {
			print("%s = %s.clone()", varStr, valStr)
		}
		print("}")

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

		if t.Kind() == reflect.Slice { // slice
			print("if %s != nil {", valStr)
			print("%s = make([]%s, len(%s))", varStr, getTypeName(t.Elem()), valStr)
		}

		if isPrimitive(t.Elem()) {
			if t.Kind() == reflect.Slice { // slice
				print("copy(%s, %s)", varStr, valStr)
			} else { // array
				print("%s = %s", varStr, valStr)
			}
		} else {
			print("for %s := 0; %s < len(%s); %s++ {", iNames[iDepth], iNames[iDepth], valStr, iNames[iDepth])
			genField(
				t.Elem(),
				fmt.Sprintf("%s[%s]", varStr, iNames[iDepth]),
				fmt.Sprintf("%s[%s]", valStr, iNames[iDepth]),
				iDepth+1, kvDepth,
			)
			print("}")
		}
		if t.Kind() == reflect.Slice { // slice
			print("}")
		}

	case reflect.Map:
		print("if %s != nil {", valStr)
		print("%s = make(map[%s]%s)", varStr, getTypeName(t.Key()), getTypeName(t.Elem()))
		print("for %s,%s := range %s {", kNames[kvDepth], vNames[kvDepth], valStr)
		genField(
			t.Elem(),
			fmt.Sprintf("%s[%s]", varStr, kNames[kvDepth]),
			fmt.Sprintf("%s", vNames[kvDepth]),
			iDepth, kvDepth+1,
		)
		print("}")
		print("}")

	default:
		if !isPrimitive(t) {
			panic("unsupported type:" + t.String())
		}

		print("%s = %s", varStr, valStr)
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

func getTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Ptr:
		return "*" + t.Elem().Name()
	case reflect.Slice:
		return fmt.Sprintf("[]%s", getTypeName(t.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", t.Key().Name(), getTypeName(t.Elem()))
	default:
		return t.Name()
	}
}
