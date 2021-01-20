package ansible

import (
	"github.com/ghodss/yaml"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"testing"
)

const validPlaybook = `
- hosts: aws-node
  gather_facts: no
  roles:
  - { role: add-aws-facts, tags: [ add-aws-facts ] }
  post_tasks:
    - setup:
  tags: always,alwaystest
`

func TestReadFromFile(t *testing.T) {

	t.Run("with a valid play and tags", func(t *testing.T) {
		binData := []byte(validPlaybook)
		var plays []Play
		err := yaml.Unmarshal([]byte(binData), &plays)
		assert.Nil(t, err)
		assert.NotNil(t, plays)
		assert.NotEmpty(t, plays)

		deep.Equal(plays[0], Play{
			Hosts: "aws-node",
			Roles: []Role{
				{Name: "add-aws-facts"},
			},
			RawTags: []string{"add-aws-facts"},
		})

	})

}
