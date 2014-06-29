package structor

import "reflect"

// isExportableField is whether a struct field is exported.
func isExportableField(field reflect.StructField) bool {
	return field.PkgPath == "" // PkgPath is empty for exported fields.
}
