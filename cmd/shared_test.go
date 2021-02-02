package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

const (
	ProjectMultiLevelPath  = "../testdata/projectMultiLevel"
	ProjectSimpleLevelPath = "../testdata/projectSimpleLevel"
)

func TestExtractMultitpleCompletion(t *testing.T) {
	testCases := map[string]struct {
		toComplete       string
		expectsRemainder string
		expectsCurrent   string
	}{"empty string should return empty strings": {
		toComplete:       "",
		expectsRemainder: "",
		expectsCurrent:   "",
	}, "'first' should return no remainder and 'first' as current": {
		toComplete:       "first",
		expectsRemainder: "",
		expectsCurrent:   "first",
	}, "'first,' should return 'first,' as remainder and empty as current": {
		toComplete:       "first,",
		expectsRemainder: "first,",
		expectsCurrent:   "",
	}, "'first,second' should return 'first,' as remainder and second as current": {
		toComplete:       "first,second",
		expectsRemainder: "first,",
		expectsCurrent:   "second",
	}, "'first,second,' should return 'first,second,' as remainder and empty as current": {
		toComplete:       "first,second,",
		expectsRemainder: "first,second,",
		expectsCurrent:   "",
	}, "'first,second,third' should return 'first,second,' as remainder and third as current": {
		toComplete:       "first,second,third",
		expectsRemainder: "first,second,",
		expectsCurrent:   "third",
	}}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			actualRemainder, actualCurrent := extractMultipleCompletion(testCase.toComplete)
			assert.Equal(t, testCase.expectsRemainder, actualRemainder)
			assert.Equal(t, testCase.expectsCurrent, actualCurrent)
		})
	}

}

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
		"multi-level with 'customer==customer1,' should return to complete": {
			toComplete: "customer==customer1",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer==customer1"},
		},
		"multi-level with 'customer==customer1,' should return all vars": {
			toComplete: "customer==customer1,",
			path:       ProjectMultiLevelPath,
			expect:     []string{"customer==customer1,customer", "customer==customer1,env", "customer==customer1,os", "customer==customer1,platform"},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {

			fmt.Println()

			actual, _ := filterCompletion(testCase.toComplete, testCase.path)
			assert.Equal(t, testCase.expect, actual)
		})
	}

}

func TestTagsCompletion(t *testing.T) {
	testCases := map[string]struct {
		toComplete string
		path       string
		playbook   string
		expect     []string
	}{"multi-level with black should return all vars": {
		toComplete: "",
		path:       ProjectMultiLevelPath,
		playbook:   "test.yml",
		expect:     []string{"existing-role", "playtag1", "role-1", "test1-tag", "test2-tag"},
	},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			playbook, _ := filepath.Abs(fmt.Sprintf("%s/%s", testCase.path, testCase.playbook))
			actual, _ := tagsCompletion(testCase.toComplete, testCase.path, playbook)
			assert.Equal(t, testCase.expect, actual)
		})
	}

}

func TestPlaybookCompletion(t *testing.T) {
	testCases := map[string]struct {
		toComplete string
		path       string
		expects    []string
	}{"multi-level with black should return all vars": {
		toComplete: "",
		path:       ProjectMultiLevelPath,
		expects:    []string{"test.yml"},
	},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {

			var playbooks []string
			for _, expect := range testCase.expects {
				playbook, _ := filepath.Abs(fmt.Sprintf("%s/%s", testCase.path, expect))
				playbooks = append(playbooks, playbook)
			}
			actual, _ := playbookCompletion(testCase.toComplete, testCase.path)
			assert.Equal(t, playbooks, actual)
		})
	}

}
