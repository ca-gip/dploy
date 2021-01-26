package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var simpleRole = Role{
	Name:  "existing-role-1",
	Tasks: nil,
	Tags:  *utils.NewSetFromSlice("role-tags-1"),
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestRole_ReadRoleTasks(t *testing.T) {

	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		err := simpleRole.LoadFromPath(utils.ProjectSimpleLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)
		utils.DeepEqual(t, []string{"test1-tag", "test2-tag"}, simpleRole.AllTags().List())
	})

	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		err := simpleRole.LoadFromPath(utils.ProjectMultiLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)

		utils.DeepEqual(t, []string{"test-tag", "test2-tag"}, simpleRole.AllTags().List())
	})
}

func TestRole_AllTags(t *testing.T) {
	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		err := simpleRole.LoadFromPath(utils.ProjectMultiLevelPath)
		assert.Nil(t, err)
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)

		utils.DeepEqual(t, []string{"test1-tag", "test2-tag"}, simpleRole.AllTags().List())
	})
}
