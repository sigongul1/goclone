package main

import "reflect"

// =================================================================================

type test_t struct {
	I             *int32
	F             float32
	Str           string
	Slice_i64     []int64
	Struct        *test2_t
	Struct2       test2_t
	Slice_test3   []test3_t
	Slice_test4   []*test4_t
	Arr_i32       [2]int32
	Arr_test3     [1]*test3_t
	T4            *test4_t
	NilTestSlice  [][]string
	NilTestSlice2 [][]*test2_t
	NilTestSlice3 [][]map[int]*test2_t
	Map           map[int32]string
	Map2          map[string]*test3_t
	Map3          map[string][]int
	NilTestMap    map[int][]map[int32]int64

	notUse uint
}

type test2_t struct {
	X  int8
	Y  float64
	T3 *test3_t
	T4 *test4_t
}

type test3_t struct {
	Hello string
	T2    *test4_t
}

type test4_t struct {
	Nice []uint16
}

// =================================================================================

func (pOwn *test_t) clone() *test_t {
	r := new(test_t)

	if pOwn.I != nil {
		r.I = new(int32)
		*r.I = *pOwn.I
	}
	r.F = pOwn.F
	r.Str = pOwn.Str
	if pOwn.Slice_i64 != nil {
		r.Slice_i64 = make([]int64, len(pOwn.Slice_i64))
		copy(r.Slice_i64, pOwn.Slice_i64)
	}
	if pOwn.Struct != nil {
		r.Struct = pOwn.Struct.clone()
	}
	r.Struct2 = *pOwn.Struct2.clone()
	if pOwn.Slice_test3 != nil {
		r.Slice_test3 = make([]test3_t, len(pOwn.Slice_test3))
		for i := 0; i < len(pOwn.Slice_test3); i++ {
			r.Slice_test3[i] = *pOwn.Slice_test3[i].clone()
		}
	}
	if pOwn.Slice_test4 != nil {
		r.Slice_test4 = make([]*test4_t, len(pOwn.Slice_test4))
		for i := 0; i < len(pOwn.Slice_test4); i++ {
			if pOwn.Slice_test4[i] != nil {
				r.Slice_test4[i] = pOwn.Slice_test4[i].clone()
			}
		}
	}
	r.Arr_i32 = pOwn.Arr_i32
	for i := 0; i < len(pOwn.Arr_test3); i++ {
		if pOwn.Arr_test3[i] != nil {
			r.Arr_test3[i] = pOwn.Arr_test3[i].clone()
		}
	}
	if pOwn.T4 != nil {
		r.T4 = pOwn.T4.clone()
	}
	if pOwn.NilTestSlice != nil {
		r.NilTestSlice = make([][]string, len(pOwn.NilTestSlice))
		for i := 0; i < len(pOwn.NilTestSlice); i++ {
			if pOwn.NilTestSlice[i] != nil {
				r.NilTestSlice[i] = make([]string, len(pOwn.NilTestSlice[i]))
				copy(r.NilTestSlice[i], pOwn.NilTestSlice[i])
			}
		}
	}
	if pOwn.NilTestSlice2 != nil {
		r.NilTestSlice2 = make([][]*test2_t, len(pOwn.NilTestSlice2))
		for i := 0; i < len(pOwn.NilTestSlice2); i++ {
			if pOwn.NilTestSlice2[i] != nil {
				r.NilTestSlice2[i] = make([]*test2_t, len(pOwn.NilTestSlice2[i]))
				for j := 0; j < len(pOwn.NilTestSlice2[i]); j++ {
					if pOwn.NilTestSlice2[i][j] != nil {
						r.NilTestSlice2[i][j] = pOwn.NilTestSlice2[i][j].clone()
					}
				}
			}
		}
	}
	if pOwn.NilTestSlice3 != nil {
		r.NilTestSlice3 = make([][]map[int]*test2_t, len(pOwn.NilTestSlice3))
		for i := 0; i < len(pOwn.NilTestSlice3); i++ {
			if pOwn.NilTestSlice3[i] != nil {
				r.NilTestSlice3[i] = make([]map[int]*test2_t, len(pOwn.NilTestSlice3[i]))
				for j := 0; j < len(pOwn.NilTestSlice3[i]); j++ {
					if pOwn.NilTestSlice3[i][j] != nil {
						r.NilTestSlice3[i][j] = make(map[int]*test2_t)
						for k, v := range pOwn.NilTestSlice3[i][j] {
							if v != nil {
								r.NilTestSlice3[i][j][k] = v.clone()
							}
						}
					}
				}
			}
		}
	}
	if pOwn.Map != nil {
		r.Map = make(map[int32]string)
		for k, v := range pOwn.Map {
			r.Map[k] = v
		}
	}
	if pOwn.Map2 != nil {
		r.Map2 = make(map[string]*test3_t)
		for k, v := range pOwn.Map2 {
			if v != nil {
				r.Map2[k] = v.clone()
			}
		}
	}
	if pOwn.Map3 != nil {
		r.Map3 = make(map[string][]int)
		for k, v := range pOwn.Map3 {
			if v != nil {
				r.Map3[k] = make([]int, len(v))
				copy(r.Map3[k], v)
			}
		}
	}
	if pOwn.NilTestMap != nil {
		r.NilTestMap = make(map[int][]map[int32]int64)
		for k, v := range pOwn.NilTestMap {
			if v != nil {
				r.NilTestMap[k] = make([]map[int32]int64, len(v))
				for i := 0; i < len(v); i++ {
					if v[i] != nil {
						r.NilTestMap[k][i] = make(map[int32]int64)
						for k2, v2 := range v[i] {
							r.NilTestMap[k][i][k2] = v2
						}
					}
				}
			}
		}
	}

	return r
}
func (pOwn *test2_t) clone() *test2_t {
	r := new(test2_t)

	r.X = pOwn.X
	r.Y = pOwn.Y
	if pOwn.T3 != nil {
		r.T3 = pOwn.T3.clone()
	}
	if pOwn.T4 != nil {
		r.T4 = pOwn.T4.clone()
	}

	return r
}
func (pOwn *test3_t) clone() *test3_t {
	r := new(test3_t)

	r.Hello = pOwn.Hello
	if pOwn.T2 != nil {
		r.T2 = pOwn.T2.clone()
	}

	return r
}
func (pOwn *test4_t) clone() *test4_t {
	r := new(test4_t)

	if pOwn.Nice != nil {
		r.Nice = make([]uint16, len(pOwn.Nice))
		copy(r.Nice, pOwn.Nice)
	}

	return r
}

// =================================================================================

func example_check() {
	x := &test_t{
		I:             new(int32),
		F:             3.14,
		Str:           "haha",
		Slice_i64:     []int64{3, 7, 8},
		Struct:        &test2_t{X: 10, Y: 20, T3: &test3_t{Hello: "hell", T2: &test4_t{Nice: []uint16{5, 15}}}, T4: &test4_t{Nice: []uint16{7, 9}}},
		Struct2:       test2_t{X: 30, Y: 110, T3: &test3_t{Hello: "oli", T2: &test4_t{Nice: []uint16{11, 22}}}, T4: &test4_t{Nice: []uint16{78, 12}}},
		Slice_test3:   []test3_t{{Hello: "aaa", T2: &test4_t{Nice: []uint16{88, 77}}}},
		Slice_test4:   []*test4_t{{Nice: []uint16{123}}},
		Arr_i32:       [2]int32{1, 2},
		Arr_test3:     [1]*test3_t{{Hello: "bb", T2: &test4_t{Nice: []uint16{55, 34}}}},
		T4:            &test4_t{Nice: []uint16{28, 92}},
		NilTestSlice:  nil,
		NilTestSlice2: nil,
		NilTestSlice3: [][]map[int]*test2_t{},
		Map:           map[int32]string{7: "ok", 8: "not"},
		Map2: map[string]*test3_t{"k1": {
			Hello: "ho",
			T2:    nil,
		}},
		Map3:       map[string][]int{"ii": {1, 2, 4}, "jj": {2, 3, 8}},
		NilTestMap: map[int][]map[int32]int64{},

		notUse: 0,
	}

	y := x.clone()

	print("%+v", x)
	print("%+v", y)
	print("%#v", reflect.DeepEqual(x, y))

}
