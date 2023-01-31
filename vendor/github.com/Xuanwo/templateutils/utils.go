package templateutils

import (
	"fmt"
	"reflect"
	"strings"
)

// Equal will check whether two value is equal.
func Equal(a, b reflect.Value) bool {
	aType := reflect.TypeOf(a)
	if aType == nil {
		return false
	}
	bType := reflect.ValueOf(b)
	if bType.IsValid() && bType.Type().ConvertibleTo(aType) {
		// Attempt comparison after type conversion
		return reflect.DeepEqual(bType.Convert(aType).Interface(), a)
	}

	return false
}

// In will check whether value is in item
func In(item reflect.Value, value reflect.Value) (bool, error) {
	// indirect item
	item = indirectInterface(item)
	if !item.IsValid() {
		return false, fmt.Errorf("value of untyped nil")
	}
	var isNil bool
	if item, isNil = indirect(item); isNil {
		return false, fmt.Errorf("value of nil pointer")
	}

	// indirect value
	value = indirectInterface(value)
	if !value.IsValid() {
		return false, fmt.Errorf("value of untyped nil")
	}
	if value, isNil = indirect(value); isNil {
		return false, fmt.Errorf("value of nil pointer")
	}

	itemType := item.Type()
	valueType := value.Type()

	switch itemType.Kind() {
	case reflect.Array, reflect.Slice:
		if itemType.Elem().Kind() != valueType.Kind() {
			return false, fmt.Errorf("not the same type, expected %s, got %s", itemType, valueType)
		}
		for i := 0; i < item.Len(); i++ {
			if reflect.DeepEqual(item.Index(i).Interface(), value.Interface()) {
				return true, nil
			}
		}
		return false, nil
	case reflect.String:
		if valueType.Kind() != reflect.String {
			return false, fmt.Errorf("type expect String, got %s", value)
		}
		return strings.Contains(item.String(), value.String()), nil
	case reflect.Map:
		if itemType.Key().Kind() != valueType.Kind() {
			return false, fmt.Errorf("not the same type, expected %s, got %s", itemType, valueType)
		}
		v := item.MapIndex(value)
		return v.IsValid(), nil
	case reflect.Invalid:
		// the loop holds invariant: item.IsValid()
		panic("unreachable")
	default:
		return false, fmt.Errorf("can't check item of type %s", item.Type())
	}
}

// MakeSlice will create a new slice via input values.
func MakeSlice(item ...interface{}) []interface{} {
	return item
}
