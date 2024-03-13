package helpers

import "reflect"

func Contains(slice interface{}, target interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	targetValue := reflect.ValueOf(target)

	if sliceValue.Kind() != reflect.Slice {
		panic("Input must be a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i)
		if element.Interface() == targetValue.Interface() {
			return true
		}
	}

	return false
}
