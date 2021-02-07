package jsonpatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/snorwin/jsonpatch"
)

var _ = Describe("JSONPointer", func() {
	Context("ParseJSONPointer_String", func() {
		It("should parse and unparse", func() {
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").String()).Should(Equal("/a/b/c"))
			Ω(jsonpatch.ParseJSONPointer("/a/b/c/").String()).Should(Equal("/a/b/c/"))
			Ω(jsonpatch.ParseJSONPointer("a/b/c").String()).Should(Equal("a/b/c"))
		})
	})
	Context("Add", func() {
		It("should add element", func() {
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Add("d").String()).Should(Equal("/a/b/c/d"))
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Add("d").Add("e").String()).Should(Equal("/a/b/c/d/e"))
		})
		It("should not modify original Pointer", func() {
			original := jsonpatch.ParseJSONPointer("/1/2/3")
			original.Add("4")
			original.Add("5")
			Ω(original.String()).Should(Equal("/1/2/3"))
		})
	})
	Context("Match", func() {
		It("should match", func() {
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("*")).Should(BeTrue())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/b/c")).Should(BeTrue())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/*/c")).Should(BeTrue())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/*")).Should(BeTrue())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/*/*")).Should(BeTrue())
		})
		It("should not match", func() {
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("*/d")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/c/b")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/b")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("/a/b/c/d/e")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("a/b/c").Match("/a/b/c")).Should(BeFalse())
			Ω(jsonpatch.ParseJSONPointer("/a/b/c").Match("a/b/c")).Should(BeFalse())
		})
	})
})
