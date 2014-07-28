package structor

import (
	"reflect"
	. "github.com/101loops/bdd"
)

var _ = Describe("Codec", func() {

	data := newSimpleStruct()
	dataType := reflect.TypeOf(SimpleStruct{})
	strType := reflect.TypeOf("test")

	Context("struct codec", func() {

		tagName := "test"

		It("from struct", func() {
			rTyp := reflect.ValueOf(data).Type()
			codec := newCodec(rTyp, tagName)

			Check(codec, NotNil)
			Check(codec.Type(), Equals, rTyp)
			Check(codec.FieldNames(), Equals, []string{"Dummy", "Yummy"})

			fields := codec.Fields()
			Check(fields, HasLen, 2)

			Check(fields[0].Name, Equals, "Dummy")
			Check(fields[0].KeyType, IsNil)
			Check(fields[0].ElemType, IsNil)

			Check(fields[1].Name, Equals, "Yummy")
			Check(fields[1].KeyType, IsNil)
			Check(fields[1].ElemType, IsNil)
		})

		It("from complex struct", func() {
			data := newComplexStruct()
			rTyp := reflect.ValueOf(data).Type()
			codec := newCodec(rTyp, tagName)

			Check(codec, NotNil)
			Check(codec.Type(), Equals, rTyp)
			Check(codec.FieldNames(), Equals, []string{"One", "Two", "Three", "Four"})

			fields := codec.Fields()
			Check(fields, HasLen, 4)

			Check(fields[2].KeyType, IsNil)

			Check(*fields[3].KeyType, Equals, strType)
			Check(*fields[3].ElemType, Equals, dataType)
		})

		It("from recursive struct", func() {
			data := newRecursiveStruct()
			rTyp := reflect.ValueOf(data).Type()
			codec := newCodec(rTyp, tagName)

			Check(codec, NotNil)
			Check(codec.Type(), Equals, rTyp)
			Check(codec.FieldNames(), Equals, []string{"Level", "Parent", "Children"})

			fields := codec.Fields()
			Check(fields, HasLen, 3)
			Check(fields[0].Name, Equals, "Level")
			Check(fields[1].Name, Equals, "Parent")
			Check(fields[2].Name, Equals, "Children")
		})
	})

	Context("create field codec", func() {

		It("from exported field", func() {
			codec := newFieldCodec(reflect.ValueOf(data).Type(), 0, "test")

			Check(codec, NotNil)
			Check(codec.Index, Equals, 0)
			Check(codec.Name, Equals, "Dummy")
			Check(codec.Label, Equals, "dummytag")
			Check(*codec.Tag, Equals, TagCodec{Name: "dummytag"})
			Check(codec.Type, Equals, reflect.ValueOf("string").Type())
		})

		It("from hidden field", func() {
			codec := newFieldCodec(reflect.ValueOf(data).Type(), 2, "test")
			Check(codec, IsNil)
		})

		It("from unexported field", func() {
			codec := newFieldCodec(reflect.ValueOf(data).Type(), 3, "test")
			Check(codec, IsNil)
		})
	})

	Context("create tag codec", func() {

		It("from empty tag", func() {
			codec := newTagCodec("")
			Check(*codec, Equals, TagCodec{})
		})

		It("from tag with name", func() {
			codec := newTagCodec("name")
			Check(*codec, Equals, TagCodec{Name: "name"})
		})

		It("from tag with modifiers only", func() {
			codec := newTagCodec(",omitempty")
			Check(*codec, Equals, TagCodec{Name: "", Modifiers: []string{"omitempty"}})
		})

		It("from tag with name and modifiers", func() {
			codec := newTagCodec("name,omitempty")
			Check(*codec, Equals, TagCodec{Name: "name", Modifiers: []string{"omitempty"}})
		})
	})
})
