package structor

import (
	"fmt"
	"reflect"
)

// Writer represents a writable struct.
type Writer struct {
	reader *Reader
}

// FieldWriter represents a writable field.
type FieldWriter struct {
	*FieldReader
	writer *Writer
}

func newWriter(dst interface{}, reader *Reader) (*Writer, error) {
	rVal := reflect.ValueOf(dst)
	if rVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("structor: writer requires pointer to struct")
	}

	return &Writer{reader}, nil
}

// Fields returns the writable fields of the struct.
func (w *Writer) Fields() []*FieldWriter {
	fields := make([]*FieldWriter, len(w.reader.codec.Fields))
	for idx, fld := range w.reader.Fields() {
		fields[idx] = &FieldWriter{fld, w}
	}
	return fields
}

// SetValue applies the passed-in value to the field. It returns an error if
// the field is not writable or the value's type does not match the field's.
func (fld *FieldWriter) SetValue(value interface{}) error {
	rField := fld.reader.rVal.Field(fld.Index)
	if !rField.CanSet() {
		return fmt.Errorf("structor: can not set field %q", fld.Name)
	}

	vVal := reflect.ValueOf(value)
	vType := vVal.Type()
	if fld.Type != vType {
		return fmt.Errorf("structor: unable to set field %q (%v) to value '%v' (%v)",
			fld.Name, fld.Type, value, vType)
	}

	rField.Set(vVal)

	return nil
}
