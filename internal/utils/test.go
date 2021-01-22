package utils

import (
	"fmt"
	"reflect"
	"testing"
)

const testPath = "./../../test"

var (
	ProjectMultiLevelPath  = fmt.Sprint(testPath, "/projectMultiLevel")
	ProjectSimpleLevelPath = fmt.Sprint(testPath, "/projectSimpleLevel")
)

func DeepEqual(t *testing.T, expected interface{}, actual interface{}) {
	if eq := reflect.DeepEqual(expected, actual); !eq {
		t.Errorf("\nStruct not equals \nExpected\t%+v\nActual\t\t%+v\n", expected, actual)
	}
}

func DeepNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if eq := reflect.DeepEqual(expected, actual); eq {
		t.Errorf("\nStruct are both \n\t%+v", expected)
	}
}
