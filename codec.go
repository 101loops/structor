package structor

import (
	"reflect"
	"strings"
)

// Codec represents a struct and its fields.
type Codec struct {
	// Type returns the struct's type.
	Type reflect.Type

	// Fields returns the struct's field codecs.
	Fields []*FieldCodec

	// FieldNames returns the struct's exportable field names.
	FieldNames []string

	// Attrs can contain custom attributes of the codec.
	Attrs map[string]interface{}

	// Complete is whether the codec is completely processed.
	// An incomplete codec may be encountered when walking a recursive struct.
	Complete bool
}

func newCodec(rType reflect.Type, tagName string) *Codec {
	ret := &Codec{
		Type:  rType,
		Attrs: make(map[string]interface{}),
	}

	fieldsCount := rType.NumField()
	for idx := 0; idx < fieldsCount; idx++ {
		fCodec := newFieldCodec(rType, idx, tagName)
		if fCodec == nil {
			continue
		}
		fCodec.Parent = ret

		ret.Fields = append(ret.Fields, fCodec)
		ret.FieldNames = append(ret.FieldNames, fCodec.Name)
	}

	return ret
}

// FieldCodec represents a struct field.
type FieldCodec struct {

	// Parent is a reference to the struct the field is in.
	Parent *Codec

	// Index is the index of the field in the struct.
	Index int

	// Anonymous is whether the field is embedded.
	Anonymous bool

	// Name is the name of the field in the struct.
	Name string

	// Tag is the codec for the field's tags.
	Tag *TagCodec

	// Attrs can contain custom attributes of the field's codec.
	Attrs map[string]interface{}

	// Type is the field's type.
	Type reflect.Type

	// KeyType is the type of the field's key, if any.
	KeyType *reflect.Type

	// ElemType is the type of the field's collection's value type, if any.
	ElemType *reflect.Type
}

func newFieldCodec(rType reflect.Type, idx int, tagName string) *FieldCodec {
	fld := rType.Field(idx)
	if !isExportableField(fld) {
		return nil
	}

	fType := fld.Type
	keyType, elemType := subTypesOf(fType)
	return &FieldCodec{
		Index:     idx,
		Anonymous: fld.Anonymous,
		Name:      fld.Name,
		Tag:       newTagCodec(fld.Tag.Get(tagName)),
		Type:      fType,
		KeyType:   keyType,
		ElemType:  elemType,
		Attrs:     make(map[string]interface{}),
	}
}

// TagCodec represents a struct field's tag.
type TagCodec struct {
	Values map[int]string
}

func newTagCodec(tag string) *TagCodec {
	vals := make(map[int]string, 0)
	if strings.TrimSpace(tag) != "" {
		for i, t := range strings.Split(tag, ",") {
			vals[i] = t
		}
	}
	return &TagCodec{vals}
}

func (tc *TagCodec) Modifiers() []string {
	var mods []string
	for i, t := range tc.Values {
		if i > 0 {
			mods = append(mods, t)
		}
	}
	return mods
}

// IndexOf returns the index of the passed-in value.
func (tc *TagCodec) IndexOf(want string) int {
	for i, tag := range tc.Values {
		if tag == want {
			return i
		}
	}
	return -1
}
