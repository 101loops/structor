package structor

import "reflect"

// Reader represents a readable struct.
type Reader struct {
	codec *Codec
	rVal  reflect.Value
}

// FieldReader represents a readable struct field.
type FieldReader struct {
	*FieldCodec
	reader *Reader
}

func newReader(src interface{}, codec *Codec) *Reader {
	rVal := reflect.ValueOf(src)
	if rVal.Kind() == reflect.Ptr {
		rVal = rVal.Elem() // TODO: remove?
	}

	return &Reader{codec, rVal}
}

// Fields returns the readable fields of the struct.
func (r *Reader) Fields() []*FieldReader {
	fields := make([]*FieldReader, len(r.codec.fields))
	for idx, fld := range r.codec.fields {
		fields[idx] = &FieldReader{fld, r}
	}
	return fields
}

// Value returns the field's value.
func (fld *FieldReader) Value() interface{} {
	rField := fld.reader.rVal.Field(fld.Index)
	return rField.Interface()
}
