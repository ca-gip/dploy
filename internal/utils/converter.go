package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func InferSlice(input interface{}) (slice []string) {
	if input == nil {
		return
	}

	value := reflect.ValueOf(input)
	rt := reflect.TypeOf(input)
	switch rt.Kind() {
	case reflect.Slice:
		slice = make([]string, value.Len())
		for i := 0; i < value.Len(); i++ {
			slice[i] = fmt.Sprintf("%v", value.Index(i).Interface())
		}
		return
	case reflect.String:
		slice = strings.Split(value.String(), ",")
		fmt.Println("String find type for", value)
		fmt.Println("String find type for", slice)
		return slice
	default:
		fmt.Println("cannot find type for", value)
	}
	return
}
