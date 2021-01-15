package ansible

import "github.com/ca-gip/dploy/internal/utils"

type Task struct {
	Role    string      `yaml:"role"`
	rawTags interface{} `yaml:"tags,omitempty"`
}

func (task *Task) Tags() []string {
	return utils.InferSlice(task.rawTags)
}
