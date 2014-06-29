package structor

import (
	. "github.com/101loops/bdd"
)

var _ = Describe("Writer", func() {

	set := newTestSet()
	data := newTestData()

	Context("set field value", func() {

		writer, err := set.NewWriter(&data)

		It("of struct pointer", func() {
			err = writer.Fields()[0].SetValue("abc")

			Check(err, IsNil)
			Check(data.Dummy, Equals, "abc")
		})

		It("invalid value type", func() {
			err = writer.Fields()[0].SetValue(123)

			Check(err, Contains, `unable to set field "Dummy" (string) to value '123' (int)`)
		})
	})

})
