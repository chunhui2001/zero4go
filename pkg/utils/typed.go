package utils

import "reflect"

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func IsStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func IsMapStringAny(t reflect.Type) bool {
	return t.Kind() == reflect.Map &&
		t.Key().Kind() == reflect.String &&
		t.Elem().Kind() == reflect.Interface
}

func IsScalar(t reflect.Type) bool {
	switch t.Kind() {
	case
		reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool:

		return true
	}

	return false
}

func IsArrayOrSlice(t reflect.Type) bool {

	return t.Kind() == reflect.Array || t.Kind() == reflect.Slice
}
