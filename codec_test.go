package structor

import (
	. "github.com/101loops/bdd"
	"reflect"
)

var _ = Describe("Codec", func() {

	data := newTestData()

	Context("struct codec", func() {
		// TODO
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
			Check(*codec, Equals, TagCodec{Name: "", Mods: []string{"omitempty"}})
		})

		It("from tag with name and modifiers", func() {
			codec := newTagCodec("name,omitempty")
			Check(*codec, Equals, TagCodec{Name: "name", Mods: []string{"omitempty"}})
		})
	})
})
