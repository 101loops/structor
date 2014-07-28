package structor

import (
	"reflect"
	"strings"
)

// Codec represents a struct and its fields.
type Codec struct {
	rType      reflect.Type
	fields     []*FieldCodec
	fieldNames []string

	// attrs can contain custom attributes of the codec.
	attrs map[string]interface{}

	// complete is whether the codec is complete.
	// An incomplete codec may be encountered when walking a recursive struct.
	complete bool
}

func newCodec(rType reflect.Type, tagName string) *Codec {
	fieldsCount := rType.NumField()

	fields := []*FieldCodec{}
	fieldNames := []string{}

	for idx := 0; idx < fieldsCount; idx++ {
		fCodec := newFieldCodec(rType, idx, tagName)
		if fCodec == nil {
			continue
		}

		fields = append(fields, fCodec)
		fieldNames = append(fieldNames, fCodec.Name)
	}

	return &Codec{
		rType:      rType,
		fields:     fields,
		fieldNames: fieldNames,
		attrs:      make(map[string]interface{}),
	}
}

// Type returns the struct's type.
func (c *Codec) Type() reflect.Type {
	return c.rType
}

// Complete is whether the codec is completely processed.
func (c *Codec) Complete() bool {
	return c.complete
}

// Fields returns the struct's field codecs.
func (c *Codec) Fields() []*FieldCodec {
	return c.fields
}

// FieldNames returns the struct's exportable field names.
func (c *Codec) FieldNames() []string {
	return c.fieldNames
}

// FieldCodec represents a struct field.
type FieldCodec struct {
	Index    int
	Name     string
	Label    string
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

	fLabel := fTag.Name
	if fLabel == "-" {
		return nil
	}
	if fLabel == "" {
		fLabel = fName
	}

	keyType, elemType := subTypesOf(fType)

	return &FieldCodec{
		Index:    idx,
		Name:     fName,
		Label:    fLabel,
		Tag:      fTag,
		Type:     fType,
		KeyType:  keyType,
		ElemType: elemType,
	}
}

// TagCodec represents a struct field's tag.
type TagCodec struct {
	Name      string
	Modifiers []string
}

func newTagCodec(tag string) *TagCodec {
	tagSplit := strings.Split(tag, ",")

	var mods []string
	if len(tagSplit) > 1 {
		mods = tagSplit[1:]
	}

	var name string
	if len(tagSplit) > 0 && tagSplit[0] != "" {
		name = tagSplit[0]
	}

	return &TagCodec{name, mods}
}
