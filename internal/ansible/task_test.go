package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

const validTask1 = `
- name: Task1
  template: src="source" dest="destination" owner=root
  tags: tasktag1
`

const validTaskWithoutName = `
- template: src="source" dest="destination" owner=root
  tags: tasktag1
`

var expectedValidTask1 = Task{
	Name: "Task1",
	Tags: *utils.NewSetFromSlice("tasktag1"),
}

func TestTask(t *testing.T) {

	t.Run("with a valid task data should be equals", func(t *testing.T) {
		binData := []byte(validTask1)
		var task []Task
		err := yaml.Unmarshal(binData, &task)
		assert.Nil(t, err)
		assert.NotNil(t, task)
		assert.NotEmpty(t, task)

		//deep.CompareUnexportedFields = false
		if diff := deep.Equal(task[0], expectedValidTask1); diff != nil {
			t.Error(diff)
		}
		assert.Equal(t, task[0], expectedValidTask1)

		// Not so deep ?
		if diff := deep.Equal(task[0].Tags, expectedValidTask1.Tags); diff != nil {
			t.Error(diff)
		}

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

		if diff := deep.Equal(left, right); len(diff) != 0 {
			t.Error(diff)
		}
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

		if diff := deep.Equal(left, right); len(diff) == 0 {
			t.Error(diff)
		}
	})

	t.Run("without name should have tags", func(t *testing.T) {
		binData := []byte(validTaskWithoutName)
		var task []Task
		err := yaml.Unmarshal(binData, &task)

		assert.Nil(t, err)
		assert.NotNil(t, task)
		assert.NotEmpty(t, task)

		expected := Task{
			Name: utils.EmptyString,
			Tags: *utils.NewSetFromSlice("tasktag1"),
		}

		if diff := deep.Equal(expected, task[0]); len(diff) != 0 {
			t.Error(diff)
		}
	})

}
