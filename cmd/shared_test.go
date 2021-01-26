package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ProjectMultiLevelPath  = "../testdata/projectMultiLevel"
	ProjectSimpleLevelPath = "../testdata/projectSimpleLevel"
)

func TestFilterCompletion(t *testing.T) {
	testCases := map[string]struct {
		toComplete string
		path       string
		expect     []string
	}{"multi-level with black should return all vars": {
		toComplete: "",
		path:       ProjectMultiLevelPath,
		expect:     []string{"customer", "env", "os", "platform"},
	},
		"multi-level with 'c' should should return customer with operators": {
			toComplete: "c",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer==", "customer!=", "customer^=", "customer~=", "customer$="},
		},
		"multi-level with 'customer' should return var with operators": {
			toComplete: "customer",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer==", "customer!=", "customer^=", "customer~=", "customer$="},
		},
		"multi-level with 'p' should return all vars": {
			toComplete: "p",
			path:       ProjectMultiLevelPath,
			expect:     []string{"platform==", "platform!=", "platform^=", "platform~=", "platform$="},
		},

		"multi-level with 'customer==' should return with values": {
			toComplete: "customer==",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer==customer1", "customer==customer2", "customer==customer3"},
		},
		"multi-level with 'customer!=' should return with values": {
			toComplete: "customer!=",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer!=customer1", "customer!=customer2", "customer!=customer3"},
		},
		"multi-level with 'customer$=' should return with values": {
			toComplete: "customer$=",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer$=customer1", "customer$=customer2", "customer$=customer3"},
		},
		"multi-level with 'customer~=' should return with values": {
			toComplete: "customer~=",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer~=customer1", "customer~=customer2", "customer~=customer3"},
		},
		"multi-level with 'customer^=' should return with values": {
			toComplete: "customer^=",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer^=customer1", "customer^=customer2", "customer^=customer3"},
		},
		"multi-level with 'customer%=' should return with nothing": {
			toComplete: "customer%=",
			path:       ProjectMultiLevelPath,
			expect:     []string{},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actual, _ := filterCompletion(testCase.toComplete, testCase.path)
			assert.Equal(t, testCase.expect, actual)
		})
	}

}

//func TestTagsCompletion(t *testing.T) {
//	testCases := map[string]struct {
//		toComplete string
//		path       string
//		playbook string
//		expect     []string
//	}{  "multi-level with black should return all vars": {
//		toComplete: "",
//		path:       ProjectMultiLevelPath,
//		playbook: "test.yml",
//		expect:     []string{"customer", "env", "os", "platform"},
//	},
//		"multi-level with 'c' should should return customer with operators": {
//			toComplete: "c",
//			path:       ProjectMultiLevelPath,
//			playbook: "test.yml",
//			expect:     []string{"customer==", "customer!=", "customer^=", "customer~=", "customer$="},
//		},
//	}
//
//	for testName, testCase := range testCases {
//		t.Run(testName, func(t *testing.T) {
//			actual, _ := tagsCompletion(testCase.toComplete, testCase.path, testCase.playbook)
//			assert.Equal(t, testCase.expect, actual)
//		})
//	}
//
//}

//func TestPlaybookCompletion(t *testing.T) {
//	testCases := map[string]struct {
//		toComplete string
//		path       string
//		expect     []string
//	}{"multi-level with black should return all vars": {
//		toComplete: "",
//		path:       ProjectMultiLevelPath,
//		expect:     []string{"customer", "env", "os", "platform"},
//	},
//		"multi-level with 'c' should should return customer with operators": {
//			toComplete: "c",
//			path:       ProjectMultiLevelPath,
//			expect:     []string{"customer==", "customer!=", "customer^=", "customer~=", "customer$="},
//		},
//		"multi-level with 'p' should return all vars": {
//			toComplete: "p",
//			path:       ProjectMultiLevelPath,
//			expect:     []string{"platform==", "platform!=", "platform^=", "platform~=", "platform$="},
//		},
//		"multi-level with 'customer' should return var with operators": {
//			toComplete: "",
//			path:       ProjectMultiLevelPath,
//			expect:     []string{"customer", "env", "os", "platform"},
//		},
//		"multi-level with 'customer==' should return var with operators": {
//			toComplete: "",
//			path:       ProjectMultiLevelPath,
//			expect:     []string{"customer", "env", "os", "platform"},
//		},
//	}
//
//	for testName, testCase := range testCases {
//		t.Run(testName, func(t *testing.T) {
//			actual, _ := playbookCompletion(testCase.toComplete, testCase.path)
//			assert.Equal(t, testCase.expect, actual)
//		})
//	}
//
//}
