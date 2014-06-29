package structor

import (
	. "github.com/101loops/bdd"
	"reflect"
)

var _ = Describe("Set", func() {

	data := newTestData()

	Context("add codec", func() {

		var set *Set
		BeforeEach(func() {
			set = NewSet("test")
		})

		It("from struct", func() {
			err := set.Add(data)

			Check(err, IsNil)
		})

		It("from struct pointer", func() {
			err := set.Add(&data)

			Check(err, IsNil)
		})

		It("from reflect type", func() {
			t := reflect.TypeOf(data)
			err := set.Add(t)

			Check(err, IsNil)
		})

		It("from string", func() {
			dummy := "abc 123"

			err := set.Add(dummy)
			Check(err, Contains, `structor: value is not a struct, struct pointer or reflect.Type, but "string"`)
		})
	})

	Context("create reader", func() {

		set := newTestSet()

		It("from struct", func() {
			_, err := set.NewReader(data)

			Check(err, IsNil)
		})

		It("from struct pointer", func() {
			_, err := set.NewReader(&data)

			Check(err, IsNil)
		})
	})

	Context("create writer", func() {

		set := newTestSet()

		It("from struct pointer", func() {
			_, err := set.NewWriter(&data)

			Check(err, IsNil)
		})

		It("from struct", func() {
			_, err := set.NewWriter(data)

			Check(err, Contains, "structor: writer requires pointer to struct")
		})
	})
})
