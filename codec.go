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
	Parent   *Codec
	Index    int
	Name     string
	Tag      *TagCodec
	Type     reflect.Type
	KeyType  *reflect.Type
	ElemType *reflect.Type
}

func newFieldCodec(rType reflect.Type, idx int, tagName string) *FieldCodec {
	fld := rType.Field(idx)
	if !isExportableField(fld) {
		return nil
	}

	fType := fld.Type
	fName := fld.Name
	fTag := newTagCodec(fld.Tag.Get(tagName))
	keyType, elemType := subTypesOf(fType)

	return &FieldCodec{
		Index:    idx,
		Name:     fName,
		Tag:      fTag,
		Type:     fType,
		KeyType:  keyType,
		ElemType: elemType,
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

// HasModifier returns whether the TagCodec contains the given modifier.
func (tc *TagCodec) HasModifier(want string) bool {
	for _, tag := range tc.Values {
		if tag == want {
			return true
		}
	}
	return false
}
