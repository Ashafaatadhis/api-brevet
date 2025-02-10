package utils

import (
	"reflect"
)

// TransformResponse untuk trasnformresponse
func TransformResponse(data interface{}) {
	v := reflect.ValueOf(data)

	// Pastikan input adalah pointer ke struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}

	// Rekursif nullify untuk semua field dalam struct
	recursiveNullify(v.Elem())
}

func recursiveNullify(v reflect.Value) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		elem := v.Elem()
		if elem.Kind() == reflect.Struct {
			idField := elem.FieldByName("ID")

			// Jika ada field ID dan bernilai 0, set pointer ke nil
			if idField.IsValid() && idField.CanInterface() && idField.Interface() == 0 {
				v.Set(reflect.Zero(v.Type()))
				return
			}

			// Loop semua field dalam struct
			for i := 0; i < elem.NumField(); i++ {
				field := elem.Field(i)
				if field.Kind() == reflect.Ptr || field.Kind() == reflect.Struct {
					recursiveNullify(field)
				}
			}
		}
	}
}
