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
		"equal should pass": {
			given:  "key==value",
			expect: []Filter{{Key: "key", Op: Equal, Value: "value"}},
		},
		"equal , with key/value containing dash should pass": {
			given:  "ke_y==val_ue",
			expect: []Filter{{Key: "ke_y", Op: Equal, Value: "val_ue"}},
		},
		"equal with ending should pass": {
			given:  "key==value,",
			expect: []Filter{{Key: "key", Op: Equal, Value: "value"}},
		},
		"two equal should pass": {
			given:  "key==value,key2==value2",
			expect: []Filter{{Key: "key", Op: Equal, Value: "value"}, {Key: "key2", Op: Equal, Value: "value2"}},
		},
		"op !! should nothing": {
			given:  "customer!!cagip",
			expect: nil,
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
