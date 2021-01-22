package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestTask(t *testing.T) {

	t.Run("with a valid task data should be equals", func(t *testing.T) {
		const validTask = `
- name: Task1
  template: src="source" dest="destination" owner=root
  tags: tasktag1
`
		expect := Task{
			Name: "Task1",
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		binData := []byte(validTask)
		var task []Task
		err := yaml.Unmarshal(binData, &task)

		assert.Nil(t, err)
		assert.NotNil(t, task)
		assert.NotEmpty(t, task)
		utils.DeepEqual(t, expect, task[0])
		assert.Equal(t, task[0], expect)
		utils.DeepEqual(t, task[0].Tags, expect.Tags)

	})

	t.Run("with two different task tags should fail", func(t *testing.T) {
		left := Task{
			Name: "Task1",
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		right := Task{
			Name: "Task1",
			Tags: *utils.NewSetFromSlice("tasktag2"),
		}

		utils.DeepNotEqual(t, left, right)
	})

	t.Run("with two different task name should fail", func(t *testing.T) {
		left := Task{
			Name: "Task1",
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		right := Task{
			Name: "Task2",
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		utils.DeepNotEqual(t, left, right)
	})

	t.Run("without name should have tags", func(t *testing.T) {

		const validTaskWithoutName = `
- template: src="source" dest="destination" owner=root
  tags: tasktag1
`
		expected := Task{
			Name: utils.EmptyString,
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		binData := []byte(validTaskWithoutName)
		var task []Task
		err := yaml.Unmarshal(binData, &task)

		assert.Nil(t, err)
		assert.NotNil(t, task)
		assert.NotEmpty(t, task)
		utils.DeepEqual(t, expected, task[0])
	})

}
