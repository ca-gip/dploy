package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/go-test/deep"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var simpleRole = Role{
	Name:  "existing-role-1",
	Tasks: nil,
	Tags:  utils.Set{},
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRole_ReadRoleTasks(t *testing.T) {

	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		err := simpleRole.LoadFromPath(projectSimpleLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)
		if diff := deep.Equal([]string{"test1-tag", "test2-tag"}, simpleRole.AllTags().List()); len(diff) != 0 {
			t.Error(diff)
		}
	})

	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		err := simpleRole.LoadFromPath(projectMultiLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)

		if diff := deep.Equal([]string{"test1-tag", "test2-tag"}, simpleRole.AllTags().List()); len(diff) != 0 {
			t.Error(diff)
		}
	})
}
