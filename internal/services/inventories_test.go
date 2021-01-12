package services

import (
	"github.com/relex/aini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetInventoryKeys(t *testing.T) {
	TestCases := map[string]struct {
		inventories Inventories
		expect      []string
	}{
		"single inventory with no var should return nothing ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{}}},
				},
			}},
			expect: nil,
		},
		"single inventory with one var should return one key ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": ""}}},
				},
			}},
			expect: []string{"key1"},
		},
		"single inventory with two var should return two key ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "", "key2": ""}}},
				},
			}},
			expect: []string{"key1", "key2"},
		},
		"two inventories with same var should return one key ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": ""}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": ""}}},
			}}},
			expect: []string{"key1"},
		},
		"two inventories with different var should return two key ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": ""}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key2": ""}}},
			}}},
			expect: []string{"key1", "key2"},
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			result := testCase.inventories.GetInventoryKeys()
			if testCase.expect == nil {
				assert.Equal(t, testCase.expect, result)
			} else {
				for _, expect := range testCase.expect {
					assert.Contains(t, result, expect)
				}
			}
		})
	}
}

func TestGetInventoryValues(t *testing.T) {
	TestCases := map[string]struct {
		inventories Inventories
		given       string
		expect      []string
	}{
		"single inventory with no var should return nothing ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{}}},
				},
			}},
			given:  "key1",
			expect: nil,
		},
		"single inventory with one var should return correct value ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}},
			given:  "key1",
			expect: []string{"value1"},
		},
		"single inventory with one var should return nothing ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}},
			given:  "key2",
			expect: nil,
		},
		"two inventories with same var and value should return one value ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
			}}},
			given:  "key1",
			expect: []string{"value1"},
		},
		"two inventories with same var but different value should return two value ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value2"}}},
			}}},
			given:  "key1",
			expect: []string{"value1", "value2"},
		},
		"two inventories with different var should return one value ": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key2": ""}}},
			}}},
			given:  "key1",
			expect: []string{"value1"},
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			result := testCase.inventories.GetInventoryValues(testCase.given)
			if testCase.expect == nil {
				assert.Equal(t, testCase.expect, result)
			} else {
				for _, expect := range testCase.expect {
					assert.Contains(t, result, expect)
				}
			}

		})
	}
}

func TestFilter(t *testing.T) {
	TestCases := map[string]struct {
		inventories Inventories
		filters     Filters
		expect      int
	}{
		"two inventories with same var and Equal op should return two": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
			}}},
			filters: Filters{{Key: "key", Op: Equal, Value: "value"}},
			expect:  2,
		},
		"two inventories with same var and NotEqual op should return zero": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
			}}},
			filters: Filters{{Key: "key", Op: NotEqual, Value: "value"}},
			expect:  0,
		},
		"two inventories different var and NotEqual op should return one": {
			inventories: Inventories{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value1"}}},
			}}},
			filters: Filters{{Key: "key", Op: NotEqual, Value: "value"}},
			expect:  1,
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			result := testCase.inventories.Filter(testCase.filters)
			assert.Equal(t, testCase.expect, len(result))

		})
	}
}
