package jsonpatch_test

import (
	"encoding/json"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bxcodec/faker/v3"
	jsonpatch2 "github.com/evanphx/json-patch"

	"github.com/snorwin/jsonpatch"
)

type A struct {
	B *B `json:"ptr,omitempty"`
	C C  `json:"struct"`
}

type B struct {
	Str     string  `json:"str,omitempty"`
	Bool    bool    `json:"bool"`
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Int16   int16   `json:"int16"`
	Int32   int32   `json:"int32"`
	Int64   int64   `json:"int64"`
	Uint    uint    `json:"uint"`
	Uint8   uint8   `json:"uint8"`
	Uint16  uint16  `json:"uint16"`
	Uint32  uint32  `json:"uint32"`
	Uint64  uint64  `json:"uint64"`
	UintPtr uintptr `json:"ptr" faker:"-"`
}

type C struct {
	Str string            `json:"str,omitempty"`
	Map map[string]string `json:"map"`
}

type D struct {
	PtrSlice           []*B     `json:"ptr"`
	StructSlice        []C      `json:"structs"`
	StringSlice        []string `json:"strs"`
	IntSlice           []int    `json:"ints"`
	StructSliceWithKey []C      `json:"structsWithKey"`
	PtrSliceWithKey    []*B     `json:"ptrWithKey"`
}

type E struct {
	unexported int
	Exported   int `json:"exported"`
}

type F struct {
	A *A `json:"a"`
	B *B `json:"b,omitempty"`
	C C  `json:"c"`
	D D  `json:"d"`
	E E  `json:"e"`
}

type G struct {
	I interface{} `json:"i"`
}

type H struct {
	Ignored    string `json:"_"`
	NotIgnored string `json:"notIgnored"`
}

var _ = Describe("JSONPatch", func() {
	Context("CreateJsonPatch_pointer_values", func() {
		It("pointer", func() {
			// add
			testPatch(A{B: &B{Str: "test"}}, A{})
			// remove
			testPatch(A{}, A{B: &B{Str: "test"}})
			// replace
			testPatch(A{B: &B{Str: "test1"}}, A{B: &B{Str: "test2"}})
			testPatch(&B{Str: "test1"}, &B{Str: "test2"})
			// no change
			testPatch(A{B: &B{Str: "test2"}}, A{B: &B{Str: "test2"}})
			testPatch(A{}, A{})
		})
	})
	Context("CreateJsonPatch_struct", func() {
		It("pointer", func() {
			// add
			testPatch(A{C: C{}}, A{})
			// remove
			testPatch(A{}, A{C: C{}})
			// replace
			testPatch(A{C: C{Str: "test1"}}, A{C: C{Str: "test2"}})
			// no change
			testPatch(A{C: C{Str: "test2"}}, A{C: C{Str: "test2"}})
		})
	})
	Context("CreateJsonPatch_data_type_values", func() {
		It("string", func() {
			// add
			testPatch(B{Str: "test"}, B{})
			// remove
			testPatch(B{}, B{Str: "test"})
			// replace
			testPatch(B{Str: "test1"}, B{Str: "test2"})
			// no change
			testPatch(B{Str: "test1"}, B{Str: "test1"})
		})
		It("bool", func() {
			// add
			testPatch(B{Bool: true}, B{})
			// remove
			testPatch(B{}, B{Bool: false})
			// replace
			testPatch(B{Bool: false}, B{Bool: true})
			// no change
			testPatch(B{Bool: false}, B{Bool: false})
		})
		It("int", func() {
			// add
			testPatch(B{Int: -1, Int8: 2, Int16: 5, Int32: -1, Int64: 12}, B{})
			// remove
			testPatch(B{}, B{Int: 1, Int8: 2, Int16: 5, Int32: 1, Int64: 12})
			// replace
			testPatch(B{Int: -1, Int8: 2, Int16: 5, Int32: 1, Int64: 12}, B{Int: 1, Int8: 1, Int16: 1, Int32: 1, Int64: 1})
			// mixed
			testPatch(B{Int: -1, Int16: 5, Int32: 1, Int64: 3}, B{Int: -1, Int8: 1, Int32: 1, Int64: 1})
			testPatch(B{Int32: 22, Int64: 22}, B{Int: 1, Int8: 1, Int32: 1, Int64: 1})
			// no change
			testPatch(B{Int: -1, Int8: 1, Int16: 1, Int32: 1, Int64: 1}, B{Int: -1, Int8: 1, Int16: 1, Int32: 1, Int64: 1})
		})
		It("uint", func() {
			// add
			testPatch(B{Uint: 1, Uint8: 2, Uint16: 5, Uint32: 1, Uint64: 12, UintPtr: 3}, B{})
			// remove
			testPatch(B{}, B{Uint: 1, Uint8: 2, Uint16: 5, Uint32: 1, Uint64: 12, UintPtr: 3})
			// replace
			testPatch(B{Uint: 1, Uint8: 2, Uint16: 5, Uint32: 1, Uint64: 12}, B{Uint: 1, Uint8: 1, Uint16: 1, Uint32: 1, Uint64: 1})
			// mixed
			testPatch(B{Uint: 1, Uint16: 5, Uint32: 1, Uint64: 3}, B{Uint: 1, Uint8: 1, Uint32: 1, Uint64: 1})
			testPatch(B{Uint32: 22, Uint64: 22}, B{Uint: 1, Uint8: 1, Uint32: 1, Uint64: 1})
			// no change
			testPatch(B{Uint: 1, Uint8: 1, Uint16: 1, Uint32: 1, Uint64: 1}, B{Uint: 1, Uint8: 1, Uint16: 1, Uint32: 1, Uint64: 1})
		})
	})
	Context("CreateJsonPatch_map", func() {
		It("map", func() {
			// add
			testPatch(C{Map: map[string]string{"key1": "value1"}}, C{})
			// remove
			testPatch(C{Map: map[string]string{}}, C{Map: map[string]string{"key1": "value1"}})
			// replace
			testPatch(C{Map: map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}}, C{Map: map[string]string{}})
			testPatch(C{Map: map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}}, C{Map: map[string]string{"key1": "value1"}})
			testPatch(C{Map: map[string]string{"key1": "value1"}}, C{Map: map[string]string{"key1": "value2"}})
			// no change
			testPatch(C{Map: map[string]string{"key1": "value1", "key2": "value2"}}, C{Map: map[string]string{"key1": "value1", "key2": "value2"}})
		})
	})
	Context("CreateJsonPatch_slice", func() {
		It("int slice", func() {
			// add
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{})
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{}})
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{1, 2}})
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{1, 3}})
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{2, 3}})
			// remove
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{1, 2, 3, 4}})
			testPatch(D{IntSlice: []int{1, 2}}, D{IntSlice: []int{1, 2, 3, 4}})
			testPatch(D{IntSlice: []int{1}}, D{IntSlice: []int{1, 2, 3, 4}})
			testPatch(D{IntSlice: []int{}}, D{IntSlice: []int{1, 2, 3, 4}})
			// replace
			testPatch(D{IntSlice: []int{3, 2, 1}}, D{IntSlice: []int{1, 2, 3}})
			// mixed
			testPatch(D{IntSlice: []int{2}}, D{IntSlice: []int{1, 2, 3, 4}})
			testPatch(D{IntSlice: []int{4, 3, 2}}, D{IntSlice: []int{1, 2, 3, 4}})
			// no change
			testPatch(D{IntSlice: []int{1, 2, 3}}, D{IntSlice: []int{1, 2, 3}})
		})
		It("int slice ignore order", func() {
			// add
			testPatchWithExpected([]int{1, 2, 3}, []int{1, 3}, []int{1, 3, 2}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]int{1, 2, 3}, []int{1, 2}, []int{1, 2, 3}, jsonpatch.IgnoreSliceOrder())
			// no change
			testPatchWithExpected([]int{3, 2, 1}, []int{1, 2, 3}, []int{1, 2, 3}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]int{1, 2, 3}, []int{3, 2, 1}, []int{3, 2, 1}, jsonpatch.IgnoreSliceOrder())
			// remove
			testPatchWithExpected([]int{3, 1}, []int{1, 2, 3}, []int{1, 3}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]int{3, 2}, []int{1, 2, 3}, []int{2, 3}, jsonpatch.IgnoreSliceOrder())
		})
		It("uint slice ignore order", func() {
			// add
			testPatchWithExpected([]uint{1, 2, 3}, []uint{1, 3}, []uint{1, 3, 2}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]uint16{1, 2, 3}, []uint16{1, 2}, []uint16{1, 2, 3}, jsonpatch.IgnoreSliceOrder())
			// remove
			testPatchWithExpected([]uint32{3, 1}, []uint32{1, 2, 3}, []uint32{1, 3}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]uint64{3, 2}, []uint64{1, 2, 3}, []uint64{2, 3}, jsonpatch.IgnoreSliceOrder())
		})
		It("bool slice ignore order", func() {
			// add
			testPatchWithExpected([]bool{true, false}, []bool{false}, []bool{false, true}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]bool{true, false}, []bool{true}, []bool{true, false}, jsonpatch.IgnoreSliceOrder())
			// remove
			testPatchWithExpected([]bool{true}, []bool{false, true}, []bool{true}, jsonpatch.IgnoreSliceOrder())
			testPatchWithExpected([]bool{true}, []bool{true, false}, []bool{true}, jsonpatch.IgnoreSliceOrder())
		})
		It("ptr slice ignore order", func() {
			// add
			testPatchWithExpected(D{PtrSliceWithKey: []*B{{Str: "key1"}}}, D{}, D{PtrSliceWithKey: []*B{{Str: "key1"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/ptrWithKey", "str"}}))
			testPatchWithExpected(D{PtrSliceWithKey: []*B{{Str: "key1"}}}, D{PtrSliceWithKey: []*B{}}, D{PtrSliceWithKey: []*B{{Str: "key1"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/ptrWithKey", "str"}}))
			testPatchWithExpected(D{PtrSliceWithKey: []*B{{Str: "key1"}, {Str: "new"}, {Str: "key3"}}}, D{PtrSliceWithKey: []*B{{Str: "key1"}, {Str: "key3"}}}, D{PtrSliceWithKey: []*B{{Str: "key1"}, {Str: "key3"}, {Str: "new"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/ptrWithKey", "str"}}))
		})
		It("struct slice ignore order", func() {
			// add
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}}}, D{}, D{StructSliceWithKey: []C{{Str: "key1"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}}}, D{StructSliceWithKey: []C{}}, D{StructSliceWithKey: []C{{Str: "key1"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "new"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}, {Str: "new"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			// remove
			testPatchWithExpected(D{StructSliceWithKey: []C{}}, D{StructSliceWithKey: []C{{Str: "key1"}}}, D{StructSliceWithKey: []C{}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key2"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key2"}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}, {Str: "key4"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key3"}, {Str: "key2"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}, {Str: "key4"}}}, D{StructSliceWithKey: []C{{Str: "key2"}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			// replace
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key": "value1"}}}}, D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key": "value2"}}}}, D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key": "value1"}}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key1": "value"}}}}, D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key1": "value"}}}}, D{StructSliceWithKey: []C{{Str: "key", Map: map[string]string{"key1": "value"}}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			// mixed
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "new"}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}, {Str: "key4"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key3"}, {Str: "new"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key3"}, {Str: "key2"}, {Str: "new"}}}, D{StructSliceWithKey: []C{{Str: "key1"}, {Str: "key2"}, {Str: "key3"}, {Str: "key4"}}}, D{StructSliceWithKey: []C{{Str: "key2"}, {Str: "key3"}, {Str: "new"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			// no change
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key3"}, {Str: "key2", Map: map[string]string{"key": "value"}}}}, D{StructSliceWithKey: []C{{Str: "key2", Map: map[string]string{"key": "value"}}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key2", Map: map[string]string{"key": "value"}}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
			testPatchWithExpected(D{StructSliceWithKey: []C{{Str: "key2", Map: map[string]string{"key": "value"}}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key2", Map: map[string]string{"key": "value"}}, {Str: "key3"}}}, D{StructSliceWithKey: []C{{Str: "key2", Map: map[string]string{"key": "value"}}, {Str: "key3"}}}, jsonpatch.IgnoreSliceOrderWithPattern([]jsonpatch.IgnorePattern{{"/structsWithKey", "str"}}))
		})
	})
	Context("CreateJsonPatch_interface", func() {
		It("int", func() {
			// replace
			testPatch(G{2}, G{3})
			// no change
			testPatch(G{2}, G{2})
		})
		It("string", func() {
			// replace
			testPatch(G{"value1"}, G{"value2"})
			// no change
			testPatch(G{"value1"}, G{"value1"})
		})
	})
	Context("CreateJsonPatch_ignore", func() {
		It("unexported", func() {
			// add
			testPatch(E{unexported: 1, Exported: 2}, E{unexported: 2})
			// replace
			testPatch(E{unexported: 2, Exported: 2}, E{unexported: 1, Exported: 1})
			// remove
			testPatch(E{unexported: 1}, E{unexported: 1, Exported: 2})
			testPatch(E{unexported: 1}, E{Exported: 2})
			// no change
			testPatch(E{unexported: 2}, E{})
			testPatch(E{unexported: 1, Exported: 2}, E{unexported: 2, Exported: 2})
		})
		It("ignored", func() {
			// no change
			testPatchWithExpected(H{Ignored: "new", NotIgnored: "new"}, H{Ignored: "old", NotIgnored: "old"}, H{Ignored: "old", NotIgnored: "new"})
		})
	})
	Context("CreateJsonPatch_with_predicates", func() {
		var (
			predicate jsonpatch.Predicate
		)
		BeforeEach(func() {
			predicate = jsonpatch.Funcs{
				AddFunc: func(path jsonpatch.JSONPointer, modified interface{}) bool {
					if b, ok := modified.(B); ok {
						return b.Bool || b.Int > 2
					}

					return true
				},
				ReplaceFunc: func(path jsonpatch.JSONPointer, modified, current interface{}) bool {
					if modifiedC, ok := modified.(C); ok {
						if currentC, ok := current.(C); ok {
							return len(modifiedC.Map) > len(currentC.Map)
						}
					}

					return true
				},
				RemoveFunc: func(path jsonpatch.JSONPointer, current interface{}) bool {
					if b, ok := current.(B); ok {
						return b.Str != "don't remove me"
					}

					return true
				},
			}
		})
		It("predicate_add", func() {
			// add
			testPatchWithExpected(F{B: &B{Bool: true, Str: "str"}}, F{}, F{B: &B{Bool: true, Str: "str"}}, jsonpatch.WithPredicate(predicate))
			testPatchWithExpected(F{B: &B{Int: 7, Str: "str"}}, F{}, F{B: &B{Int: 7, Str: "str"}}, jsonpatch.WithPredicate(predicate))
			// don't add
			testPatchWithExpected(F{B: &B{Bool: false, Str: "str"}, C: C{Map: map[string]string{"key": "value"}}}, F{}, F{C: C{Map: map[string]string{"key": "value"}}}, jsonpatch.WithPredicate(predicate))
			testPatchWithExpected(F{B: &B{Int: 0, Str: "str"}, C: C{Map: map[string]string{"key": "value"}}}, F{}, F{C: C{Map: map[string]string{"key": "value"}}}, jsonpatch.WithPredicate(predicate))
		})
		It("predicate_replace", func() {
			// replace
			testPatchWithExpected(F{C: C{Str: "new", Map: map[string]string{"key": "value"}}}, F{C: C{Str: "old"}}, F{C: C{Str: "new", Map: map[string]string{"key": "value"}}}, jsonpatch.WithPredicate(predicate))
			// don't replace
			testPatchWithExpected(F{C: C{Str: "new"}}, F{C: C{Str: "old", Map: map[string]string{"key": "value"}}}, F{C: C{Str: "old", Map: map[string]string{"key": "value"}}}, jsonpatch.WithPredicate(predicate))
		})
		It("predicate_remove", func() {
			// remove
			testPatchWithExpected(F{}, F{B: &B{Str: "remove me"}}, F{B: nil}, jsonpatch.WithPredicate(predicate))
			// don't remove
			testPatchWithExpected(F{}, F{B: &B{Str: "don't remove me"}}, F{B: &B{Str: "don't remove me"}}, jsonpatch.WithPredicate(predicate))
		})
	})
	Context("CreateJsonPatch_with_prefix", func() {
		It("empty prefix", func() {
			testPatchWithExpected(F{B: &B{Bool: true, Str: "str"}}, F{}, F{B: &B{Bool: true, Str: "str"}}, jsonpatch.WithPrefix([]string{""}))
		})
		It("pointer prefix", func() {
			prefix := "/a/ptr"
			modified := F{A: &A{B: &B{Bool: true, Str: "str"}}}
			current := F{A: &A{}}
			expected := F{A: &A{B: &B{Bool: true, Str: "str"}}}

			currentJSON, err := json.Marshal(current)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = json.Marshal(modified)
			Ω(err).ShouldNot(HaveOccurred())
			expectedJSON, err := json.Marshal(expected)
			Ω(err).ShouldNot(HaveOccurred())

			list, err := jsonpatch.CreateJSONPatch(modified.A.B, current.A.B, jsonpatch.WithPrefix(jsonpatch.ParseJSONPointer(prefix)))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(list.String()).ShouldNot(Equal(""))
			Ω(list.List()).Should(ContainElement(WithTransform(func(p jsonpatch.JSONPatch) string { return p.Path }, HavePrefix(prefix))))
			jsonPatch, err := jsonpatch2.DecodePatch(list.Raw())
			Ω(err).ShouldNot(HaveOccurred())
			patchedJSON, err := jsonPatch.Apply(currentJSON)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(patchedJSON).Should(MatchJSON(expectedJSON))
		})
		It("string prefix", func() {
			prefix := []string{"b"}
			modified := F{B: &B{Bool: true, Str: "str"}}
			current := F{}
			expected := F{B: &B{Bool: true, Str: "str"}}

			currentJSON, err := json.Marshal(current)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = json.Marshal(modified)
			Ω(err).ShouldNot(HaveOccurred())
			expectedJSON, err := json.Marshal(expected)
			Ω(err).ShouldNot(HaveOccurred())

			list, err := jsonpatch.CreateJSONPatch(modified.B, current.B, jsonpatch.WithPrefix(prefix))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(list.String()).ShouldNot(Equal(""))
			Ω(list.List()).Should(ContainElement(WithTransform(func(p jsonpatch.JSONPatch) string { return p.Path }, HavePrefix("/"+strings.Join(prefix, "/")))))
			jsonPatch, err := jsonpatch2.DecodePatch(list.Raw())
			Ω(err).ShouldNot(HaveOccurred())
			patchedJSON, err := jsonPatch.Apply(currentJSON)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(patchedJSON).Should(MatchJSON(expectedJSON))
		})
	})
	Context("CreateJsonPatch_errors", func() {
		It("not matching types", func() {
			_, err := jsonpatch.CreateJSONPatch(A{}, B{})
			Ω(err).Should(HaveOccurred())
		})
		It("not matching interface types", func() {
			_, err := jsonpatch.CreateJSONPatch(G{1}, G{"str"})
			Ω(err).Should(HaveOccurred())
		})
		It("invalid map (map[string]int)", func() {
			_, err := jsonpatch.CreateJSONPatch(G{map[string]int{"key": 2}}, G{map[string]int{"key": 3}})
			Ω(err).Should(HaveOccurred())
		})
		It("invalid map (map[int]string)", func() {
			_, err := jsonpatch.CreateJSONPatch(G{map[int]string{1: "value"}}, G{map[int]string{2: "value"}})
			Ω(err).Should(HaveOccurred())
		})
		It("ignore slice order failed (duplicated key)", func() {
			_, err := jsonpatch.CreateJSONPatch([]int{1, 1, 1, 1}, []int{1, 2, 3}, jsonpatch.IgnoreSliceOrder())
			Ω(err).Should(HaveOccurred())
			_, err = jsonpatch.CreateJSONPatch([]string{"1", "2", "3"}, []string{"1", "1"}, jsonpatch.IgnoreSliceOrder())
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("CreateJsonPatch_fuzzy", func() {
		var (
			current  F
			modified F
		)
		BeforeEach(func() {
			current = F{}
			err := faker.FakeData(&current)
			Ω(err).ShouldNot(HaveOccurred())

			modified = F{}
			err = faker.FakeData(&modified)
			Ω(err).ShouldNot(HaveOccurred())
		})

		for i := 0; i < 100; i++ {
			It("fuzzy "+strconv.Itoa(i), func() {
				testPatch(modified, current)
			})
		}

		Measure("fuzzy benchmark", func(b Benchmarker) {
			currentJSON, err := json.Marshal(current)
			Ω(err).ShouldNot(HaveOccurred())
			modifiedJSON, err := json.Marshal(modified)
			Ω(err).ShouldNot(HaveOccurred())

			var list jsonpatch.JSONPatchList
			_ = b.Time("runtime", func() {
				list, err = jsonpatch.CreateJSONPatch(modified, current)
			})
			Ω(err).ShouldNot(HaveOccurred())
			if list.Empty() {
				Ω(currentJSON).Should(MatchJSON(modifiedJSON))
				Ω(list.Len()).Should(Equal(0))

				return
			}

			jsonPatch, err := jsonpatch2.DecodePatch(list.Raw())
			Ω(err).ShouldNot(HaveOccurred())
			patchedJSON, err := jsonPatch.Apply(currentJSON)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(patchedJSON).Should(MatchJSON(modifiedJSON))
		}, 100)
	})
})

func testPatch(modified, current interface{}) {
	currentJSON, err := json.Marshal(current)
	Ω(err).ShouldNot(HaveOccurred())
	modifiedJSON, err := json.Marshal(modified)
	Ω(err).ShouldNot(HaveOccurred())

	list, err := jsonpatch.CreateJSONPatch(modified, current)
	Ω(err).ShouldNot(HaveOccurred())
	if list.Empty() {
		Ω(currentJSON).Should(MatchJSON(modifiedJSON))
		Ω(list.Len()).Should(Equal(0))
		Ω(list.String()).Should(Equal(""))

		return
	}

	Ω(list.String()).ShouldNot(Equal(""))
	jsonPatch, err := jsonpatch2.DecodePatch(list.Raw())
	Ω(err).ShouldNot(HaveOccurred())
	patchedJSON, err := jsonPatch.Apply(currentJSON)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(patchedJSON).Should(MatchJSON(modifiedJSON))
}

func testPatchWithExpected(modified, current, expected interface{}, options ...jsonpatch.Option) {
	currentJSON, err := json.Marshal(current)
	Ω(err).ShouldNot(HaveOccurred())
	_, err = json.Marshal(modified)
	Ω(err).ShouldNot(HaveOccurred())
	expectedJSON, err := json.Marshal(expected)
	Ω(err).ShouldNot(HaveOccurred())

	list, err := jsonpatch.CreateJSONPatch(modified, current, options...)
	Ω(err).ShouldNot(HaveOccurred())
	if list.Empty() {
		Ω(currentJSON).Should(MatchJSON(expectedJSON))
		Ω(list.Len()).Should(Equal(0))
		Ω(list.String()).Should(Equal(""))

		return
	}

	Ω(list.String()).ShouldNot(Equal(""))
	jsonPatch, err := jsonpatch2.DecodePatch(list.Raw())
	Ω(err).ShouldNot(HaveOccurred())
	patchedJSON, err := jsonPatch.Apply(currentJSON)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(patchedJSON).Should(MatchJSON(expectedJSON))
}
