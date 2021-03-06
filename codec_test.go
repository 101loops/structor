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
			Check(codec.Type, Equals, rTyp)
			Check(codec.FieldNames, Equals, []string{"Dummy", "Yummy"})

			fields := codec.Fields
			Check(fields, HasLen, 2)

			Check(fields[0].Name, Equals, "Dummy")
			Check(fields[0].KeyType, IsNil)
			Check(fields[0].Anonymous, IsFalse)
			Check(fields[0].ElemType, IsNil)
			Check(fields[0].Tag, Equals, &TagCodec{Values: map[int]string{0: "dummytag"}})
			Check(fields[0].Parent, Equals, codec)

			Check(fields[1].Name, Equals, "Yummy")
			Check(fields[1].KeyType, IsNil)
			Check(fields[1].ElemType, IsNil)
			Check(fields[1].Parent, Equals, codec)
		})

		It("from complex struct", func() {
			data := newComplexStruct()
			rTyp := reflect.ValueOf(data).Type()
			codec := newCodec(rTyp, tagName)

			Check(codec, NotNil)
			Check(codec.Type, Equals, rTyp)
			Check(codec.FieldNames, Equals, []string{"SimpleStruct", "One", "Two", "Three", "Four"})

			fields := codec.Fields
			Check(fields, HasLen, 5)

			Check(fields[0].Anonymous, IsTrue)
			Check(fields[3].KeyType, IsNil)
			Check(*fields[4].KeyType, Equals, strType)
			Check(*fields[4].ElemType, Equals, dataType)
		})

		It("from recursive struct", func() {
			data := newRecursiveStruct()
			rTyp := reflect.ValueOf(data).Type()
			codec := newCodec(rTyp, tagName)

			Check(codec, NotNil)
			Check(codec.Type, Equals, rTyp)
			Check(codec.FieldNames, Equals, []string{"Level", "Parent", "Children"})

			fields := codec.Fields
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
			Check(*codec.Tag, Equals, TagCodec{Values: map[int]string{0: "dummytag"}})
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
			Check(*codec, Equals, TagCodec{map[int]string{}})
			Check(codec.IndexOf("not-found"), EqualsNum, -1)
			Check(codec.Modifiers(), IsEmpty)
		})

		It("from tag with name", func() {
			codec := newTagCodec("name")
			Check(*codec, Equals, TagCodec{map[int]string{0: "name"}})
			Check(codec.Modifiers(), IsEmpty)
		})

		It("from tag with modifiers only", func() {
			codec := newTagCodec(",omitempty")
			Check(*codec, Equals, TagCodec{map[int]string{0: "", 1: "omitempty"}})
			Check(codec.IndexOf("omitempty"), EqualsNum, 1)
			Check(codec.Modifiers(), Equals, []string{"omitempty"})
		})

		It("from tag with name and modifiers", func() {
			codec := newTagCodec("name,omitempty")
			Check(*codec, Equals, TagCodec{map[int]string{0: "name", 1: "omitempty"}})
			Check(codec.Modifiers(), Equals, []string{"omitempty"})
		})
	})
})
