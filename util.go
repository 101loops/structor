package structor

import "reflect"

// isExportableField is whether a struct field is exported.
func isExportableField(field reflect.StructField) bool {
	return field.PkgPath == "" // PkgPath is empty for exported fields.
}

func subTypesOf(typ reflect.Type) (*reflect.Type, *reflect.Type) {
	switch typ.Kind() {
	case reflect.Map:
		keyType := typ.Key()
		elemType := typ.Elem()
		return &keyType, &elemType
	case reflect.Slice, reflect.Array:
		elemType := typ.Elem()
		return nil, &elemType
	}

	return nil, nil
}
