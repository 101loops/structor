package structor

import (
	"fmt"
	"reflect"
	. "github.com/101loops/bdd"
)

var _ = Describe("Set", func() {

	data := newSimpleStruct()

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
			dummy := "test"

			err := set.Add(dummy)
			Check(err, Contains, `structor: value is not a struct, struct pointer or reflect.Type - but "string"`)
		})

		It("must add", func() {
			Check(func() {
				set.AddMust("test")
			}, Panics)
		})

		It("with validator", func() {
			set.SetValidateFunc(func(s *Set, c *Codec) error {
				return fmt.Errorf("validation error")
			})

			Check(set.Add(&data), Contains, `validation error`)
		})

		It("recursive type", func() {
			t := reflect.TypeOf(newComplexStruct())
			err := set.Add(t)

			Check(err, IsNil)
			Check(set.codecs, HasLen, 2)
		})

		It("sub-type with validation error", func() {
			type InvalidSubStruct struct {
				Field1 string
				Field2 string
			}
			type InvalidStruct struct {
				Sub InvalidSubStruct
			}

			set.SetValidateFunc(func(s *Set, c *Codec) error {
				if len(c.FieldNames()) > 1 {
					return fmt.Errorf("validation error")
				}
				return nil
			})

			err := set.Add(InvalidStruct{})
			Check(err, Contains, `validation error`)
			Check(set.codecs, HasLen, 0)
		})
	})

	Context("get codec", func() {

		set := newTestSet()

		It("for struct", func() {
			codec, err := set.Get(data)

			Check(err, IsNil)
			Check(codec.Complete(), IsTrue)
		})

		It("for invalid type", func() {
			_, err := set.Get("test")
			Check(err, Contains, `structor: value is not a struct, struct pointer or reflect.Type - but "string"`)
		})

		It("for non-existent type", func() {
			type SomeStruct struct{}
			_, err := set.Get(SomeStruct{})
			Check(err, Contains, `structor: no registered codec found for type 'structor.SomeStruct'`)
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

		It("from string", func() {
			_, err := set.NewReader("test")
			Check(err, Contains, `structor: value is not a struct, struct pointer or reflect.Type - but "string"`)
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

		It("from string", func() {
			_, err := set.NewWriter("test")
			Check(err, Contains, `structor: value is not a struct, struct pointer or reflect.Type - but "string"`)
		})
	})
})
