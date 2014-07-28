package structor

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Reader", func() {

	set := newTestSet()
	data := newSimpleStruct()

	Context("get field value", func() {

		It("of struct", func() {
			reader, err := set.NewReader(data)
			Check(err, IsNil)

			val := reader.Fields()[0].Value()
			Check(val, Equals, "test")
		})

		It("of struct pointer", func() {
			reader, err := set.NewReader(&data)
			Check(err, IsNil)

			val := reader.Fields()[0].Value()
			Check(val, Equals, "test")
		})
	})
})
