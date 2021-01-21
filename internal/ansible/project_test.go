package ansible

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestProjects_LoadFromPath(t *testing.T) {

	t.Run("with a simple ansible project should have all assets fetched", func(t *testing.T) {
		project := Projects.LoadFromPath(projectSimpleLevelPath)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.Playbooks)
		assert.NotEmpty(t, project.Inventories)
		assert.Len(t, project.Playbooks, 1)
		assert.Len(t, project.Inventories, 1)
		assert.Equal(t, projectSimpleLevelPath, *project.Path)
	})

	t.Run("with a simple ansible project should have the correct playbook paths", func(t *testing.T) {
		project := Projects.LoadFromPath(projectSimpleLevelPath)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.Playbooks)
		assert.Len(t, project.PlaybookPaths(), 1)
		path, err := filepath.Abs(fmt.Sprint(projectSimpleLevelPath, "/test.yml"))
		assert.Nil(t, err)
		assert.Equal(t, path, project.PlaybookPaths()[0])
	})
}

func TestProject_PlaybookPaths(t *testing.T) {

	t.Run("with a simple ansible project should have the correct playbook paths", func(t *testing.T) {
		project := Projects.LoadFromPath(projectSimpleLevelPath)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.Playbooks)
		assert.Len(t, project.PlaybookPaths(), 1)
		path, err := filepath.Abs(fmt.Sprint(projectSimpleLevelPath, "/test.yml"))
		assert.Nil(t, err)
		assert.Equal(t, path, project.PlaybookPaths()[0])
	})
}

func TestProject_PlaybookPath(t *testing.T) {

	t.Run("should have the correct playbook path for a simple project", func(t *testing.T) {
		project := Projects.LoadFromPath(projectSimpleLevelPath)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.Playbooks)
		assert.Len(t, project.PlaybookPaths(), 1)

		path, err := filepath.Abs(fmt.Sprint(projectSimpleLevelPath, "/test.yml"))
		assert.Nil(t, err)

		actual, err := project.PlaybookPath(path)
		assert.Nil(t, err)
		assert.NotNil(t, actual)

		if diff := deep.Equal(project.Playbooks[0], *actual); len(diff) != 0 {
			t.Error(diff)
		}
	})

	t.Run("should return err if playbook doesn't exist", func(t *testing.T) {
		project := Projects.LoadFromPath(projectSimpleLevelPath)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.Playbooks)
		assert.Len(t, project.PlaybookPaths(), 1)

		path, err := filepath.Abs(fmt.Sprint(projectSimpleLevelPath, "/unexistingplaybook.yml"))
		assert.Nil(t, err)

		actual, err := project.PlaybookPath(path)
		assert.Nil(t, actual)
		assert.NotNil(t, err)

	})
}
