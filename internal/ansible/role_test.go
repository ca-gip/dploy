package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var simpleRole = Role{
	Name:  "existing-role-1",
	Tasks: nil,
	Tags:  utils.Set{},
}

func TestRole_ReadRoleTasks(t *testing.T) {

	t.Run("with a valid path and simple data should have good tags", func(t *testing.T) {
		simpleRole.ReadRoleTasks("./../../test/projectMultiLevel")
		assert.NotNil(t, simpleRole.Tasks)
		assert.NotEmpty(t, simpleRole.Tasks)
	})
}
