package ansible

import (
	"fmt"
	"github.com/ca-gip/dploy/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

const validPlaybook1 = `
- hosts: aws-node
  gather_facts: no
  roles:
  - { role: add-aws-facts, tags: [ add-aws-facts ] }
  post_tasks:
    - setup:
  tags: always,alwaystest
`

func init() {
	log.SetLevel(log.DebugLevel)
}

var expectedValidPlaybook1 = Play{
	Hosts: "aws-node",
	Roles: []*Role{
		{
			Name: "add-aws-facts",
			Tags: *utils.NewSetFromSlice("add-aws-facts"),
		},
	},
	Tags: *utils.NewSetFromSlice("always", "alwaystest"),
}

func Test(t *testing.T) {

	t.Run("with a valid play data", func(t *testing.T) {
		binData := []byte(validPlaybook1)
		var plays []Play
		err := yaml.Unmarshal(binData, &plays)
		assert.Nil(t, err)
		assert.NotNil(t, plays)
		assert.NotEmpty(t, plays)

		//deep.CompareUnexportedFields = false
		utils.DeepEqual(t, plays[0], expectedValidPlaybook1)
		utils.DeepEqual(t, plays[0].Roles[0].Tags.List(), expectedValidPlaybook1.Roles[0].Tags.List())

	})

	t.Run("with two different play should be different deep.Equals", func(t *testing.T) {
		//deep.CompareUnexportedFields = false
		var left = Play{
			Hosts: "aws-node",
			Roles: []*Role{
				{
					Name: "add-aws-facts",
					Tags: *utils.NewSetFromSlice("left"),
				},
			},
			Tags: *utils.NewSetFromSlice("left", "Left"),
		}

		var right = Play{
			Hosts: "aws-node",
			Roles: []*Role{
				{
					Name: "add-aws-facts",
					Tags: *utils.NewSetFromSlice("left"),
				},
			},
			Tags: *utils.NewSetFromSlice("right", "Right"),
		}
		utils.NotDeepEqual(t, left, right)
	})

	t.Run("with a task and task tag should return all task ( play, role, task's tags)", func(t *testing.T) {

		playbooks, err := Playbooks.LoadFromPath(utils.ProjectMultiLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, playbooks)
		assert.Len(t, playbooks, 1)
		assert.NotEmpty(t, playbooks)
		assert.NotEmpty(t, playbooks[0].AllTags())
		assert.Contains(t, playbooks[0].AllTags().List(), "playtag1")
		assert.Contains(t, playbooks[0].AllTags().List(), "test1-tag")
	})

	t.Run("unmarshalling with a non existing playbook should return error", func(t *testing.T) {
		playbook, err := Playbooks.unmarshallFromPath("nonexistingplaybook.yml", utils.ProjectMultiLevelPath)
		assert.Nil(t, playbook)
		assert.NotNil(t, err)
	})

	t.Run("unmarshalling with a non yml file should return error", func(t *testing.T) {
		inventoryFilepath := fmt.Sprint(utils.ProjectMultiLevelPath, "/inventories/openstack/customer1/hosts.ini")
		playbook, err := Playbooks.unmarshallFromPath(inventoryFilepath, utils.ProjectMultiLevelPath)
		assert.Nil(t, playbook)
		assert.NotNil(t, err)
	})

	t.Run("unmarshalling with a non playbook file should return error", func(t *testing.T) {
		roleFilepath := fmt.Sprint(utils.ProjectMultiLevelPath, "/roles/existing-role-1/tasks/main.yml")
		playbook, err := Playbooks.unmarshallFromPath(roleFilepath, utils.ProjectMultiLevelPath)
		assert.Nil(t, playbook)
		assert.NotNil(t, err)
	})

	t.Run("unmarshalling with a non valid path should return error", func(t *testing.T) {
		roleFilepath := fmt.Sprint(utils.ProjectMultiLevelPath, "../../../../../../../../main.yml")
		playbook, err := Playbooks.unmarshallFromPath(roleFilepath, utils.ProjectMultiLevelPath)
		assert.Nil(t, playbook)
		assert.NotNil(t, err)
	})

}
