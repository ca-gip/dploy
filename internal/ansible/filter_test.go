package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
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

func TestEval(t *testing.T) {

	TestCases := map[string]struct {
		given   Filter
		against string
		expect  bool
	}{
		"foo equal foo should pass": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    Equal,
				Value: "foo",
			},
			against: "foo",
			expect:  true,
		},
		"foo equal bar should fail": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    Equal,
				Value: "foo",
			},
			against: "bar",
			expect:  false,
		},
		"foo notequal bar should pass": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    NotEqual,
				Value: "foo",
			},
			against: "bar",
			expect:  true,
		},
		"foo notequal foo should fail": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    NotEqual,
				Value: "foo",
			},
			against: "foo",
			expect:  false,
		},
		"foo endwith oo should pass": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    EndWith,
				Value: "foo",
			},
			against: "oo",
			expect:  true,
		},
		"foo endwith fo should fail": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    EndWith,
				Value: "foo",
			},
			against: "fo",
			expect:  false,
		},
		"foo contains o should pass": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    Contains,
				Value: "foo",
			},
			against: "o",
			expect:  true,
		},
		"foo contains z should fail": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    Contains,
				Value: "foo",
			},
			against: "z",
			expect:  false,
		},
		"foo StartWith fo should pass": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    StartWith,
				Value: "foo",
			},
			against: "fo",
			expect:  true,
		},
		"foo StartWith oo should fail": {
			given: Filter{
				Key:   utils.EmptyString,
				Op:    StartWith,
				Value: "foo",
			},
			against: "oo",
			expect:  false,
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			result := testCase.given.Eval(testCase.against)
			assert.Equal(t, testCase.expect, result)
		})
	}
}
