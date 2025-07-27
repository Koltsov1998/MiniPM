package config

import (
	"fmt"
	"reflect"
)

func MergeObjects(base, override interface{}) error {
	baseValue := reflect.ValueOf(base)
	overrideValue := reflect.ValueOf(override)

	if baseValue.Kind() != reflect.Ptr || overrideValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both arguments must be pointers")
	}

	baseElem := baseValue.Elem()
	overrideElem := overrideValue.Elem()

	if baseElem.Kind() != reflect.Struct || overrideElem.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be pointers to structs")
	}

	if baseElem.Type() != overrideElem.Type() {
		return fmt.Errorf("both arguments must be of the same type, got %s and %s",
			baseElem.Type().Name(), overrideElem.Type().Name())
	}

	for i := 0; i < baseElem.NumField(); i++ {
		mergeField(baseElem.Field(i), overrideElem.Field(i))
	}

	return nil
}

func mergeField(base, override reflect.Value) {

	if !base.CanSet() {
		return
	}

	switch override.Kind() {
	case reflect.Struct:

		for i := 0; i < override.NumField(); i++ {
			mergeField(base.Field(i), override.Field(i))
		}

	case reflect.Map:

		if override.IsNil() {
			return
		}

		if base.IsNil() {
			base.Set(reflect.MakeMap(override.Type()))
		}

		for _, key := range override.MapKeys() {
			overrideVal := override.MapIndex(key)
			base.SetMapIndex(key, overrideVal)
		}

	case reflect.Slice, reflect.Array:

		if override.Len() > 0 {
			base.Set(override)
		}

	case reflect.Ptr, reflect.Interface:

		if override.IsNil() {
			return
		}

		if base.IsNil() {

			newValue := reflect.New(override.Elem().Type())
			base.Set(newValue)
		}

		mergeField(base.Elem(), override.Elem())

	default:

		if !isZeroValue(override) {
			base.Set(override)
		}
	}
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == complex(0, 0)
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Chan:
		return v.IsNil() || v.Len() == 0
	case reflect.Func:
		return v.IsNil()
	case reflect.Array:
		return isEmptyArray(v)
	case reflect.Struct:
		return isEmptyStruct(v)
	default:
		return false
	}
}

func isEmptyArray(v reflect.Value) bool {
	for i := 0; i < v.Len(); i++ {
		if !isZeroValue(v.Index(i)) {
			return false
		}
	}
	return true
}

func isEmptyStruct(v reflect.Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if !isZeroValue(v.Field(i)) {
			return false
		}
	}
	return true
}
