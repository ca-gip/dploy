package services

import (
	"github.com/go-test/deep"
	"testing"
)

func TestParseFilterArgsFromString(t *testing.T) {

	TestCases := map[string]struct {
		given  string
		expect []Filter
	}{
		"customer==cagip should pass": {
			given:  "customer==cagip",
			expect: []Filter{{Key: "customer", Op: Equal, Value: "cagip"}},
		},
		"customer==cagip,platform==os should pass": {
			given:  "customer==cagip,platform==os",
			expect: []Filter{{Key: "customer", Op: Equal, Value: "cagip"}, {Key: "platform", Op: Equal, Value: "os"}},
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			result := ParseFilterArgsFromString(testCase.given)
			if diff := deep.Equal(testCase.expect, result); diff != nil {
				t.Error(diff)
			}
		})
	}

}
