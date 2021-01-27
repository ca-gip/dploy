package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/relex/aini"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProject_InventoryKeys(t *testing.T) {

	TestCases := map[string]struct {
		inventories []*Inventory
		expect      []string
	}{
		"single inventory with no var should return nothing ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{}}},
				},
			}},
			expect: nil,
		},
		"single inventory with one var should return one key ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": utils.EmptyString}}},
				},
			}},
			expect: []string{"key1"},
		},
		"single inventory with two var should return two key ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": utils.EmptyString, "key2": utils.EmptyString}}},
				},
			}},
			expect: []string{"key1", "key2"},
		},
		"two inventories with same var should return one key ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": utils.EmptyString}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": utils.EmptyString}}},
			}}},
			expect: []string{"key1"},
		},
		"two inventories with different var should return two key ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": utils.EmptyString}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key2": utils.EmptyString}}},
			}}},
			expect: []string{"key1", "key2"},
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			partialProject := Project{
				Inventories: testCase.inventories,
			}
			result := partialProject.InventoryKeys()
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

func TestProject_InventoryValues(t *testing.T) {
	TestCases := map[string]struct {
		inventories []*Inventory
		given       string
		expect      []string
	}{
		"single inventory with no var should return nothing ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{}}},
				},
			}},
			given:  "key1",
			expect: nil,
		},
		"single inventory with one var should return correct value ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}},
			given:  "key1",
			expect: []string{"value1"},
		},
		"single inventory with one var should return nothing ": {
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}},
			given:  "key2",
			expect: nil,
		},
		"two inventories with same var and value should return one value ": {
			inventories: []*Inventory{{
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
			inventories: []*Inventory{{
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
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key1": "value1"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key2": utils.EmptyString}}},
			}}},
			given:  "key1",
			expect: []string{"value1"},
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			partialProject := Project{
				Inventories: testCase.inventories,
			}
			result := partialProject.InventoryValues(testCase.given)
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

func TestProject_FilterInventory(t *testing.T) {
	TestCases := map[string]struct {
		inventories []*Inventory
		filters     Filters
		expect      int
	}{
		"two inventories with same var and Equal op should return two": {
			inventories: []*Inventory{{
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
			inventories: []*Inventory{{
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
			inventories: []*Inventory{{
				Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value"}}},
				},
			}, {Data: &aini.InventoryData{
				Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value1"}}},
			}}},
			filters: Filters{{Key: "key", Op: NotEqual, Value: "value"}},
			expect:  1,
		},
		"key==value,foo==bar,bar==foo should return two inventory": {
			inventories: []*Inventory{
				{Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "foo"}}},
				}}},
			filters: Filters{{Key: "key", Op: Equal, Value: "value"}, {Key: "foo", Op: Equal, Value: "bar"}, {Key: "bar", Op: Equal, Value: "foo"}},
			expect:  2,
		},
		"key==value,foo==bar,bar!=foo should return two inventory": {
			inventories: []*Inventory{
				{Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "foo"}}},
				}}},
			filters: Filters{{Key: "key", Op: Equal, Value: "value"}, {Key: "foo", Op: Equal, Value: "bar"}, {Key: "bar", Op: NotEqual, Value: "foo"}},
			expect:  1,
		},
		"key==value,foo==bar,bar==oof should return one inventory": {
			inventories: []*Inventory{
				{Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "foo"}}},
				}}},
			filters: Filters{{Key: "key", Op: Equal, Value: "value"}, {Key: "foo", Op: Equal, Value: "bar"}, {Key: "bar", Op: Equal, Value: "oof"}},
			expect:  1,
		},
		"key==value,foo!=bar,bar==oof should return one inventory": {
			inventories: []*Inventory{
				{Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "bar", "bar": "foo"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "oof"}}},
				}}, {Data: &aini.InventoryData{
					Groups: map[string]*aini.Group{"all": {Vars: map[string]string{"key": "value", "foo": "rab", "bar": "foo"}}},
				}}},
			filters: Filters{{Key: "key", Op: Equal, Value: "value"}, {Key: "foo", Op: NotEqual, Value: "bar"}, {Key: "bar", Op: Equal, Value: "oof"}},
			expect:  1,
		},
	}

	for testName, testCase := range TestCases {
		t.Run(testName, func(t *testing.T) {
			partialProject := Project{
				Inventories: testCase.inventories,
			}
			result := partialProject.FilterInventory(testCase.filters)
			assert.Equal(t, testCase.expect, len(result))

		})
	}
}
