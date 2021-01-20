package ansible

import "github.com/ca-gip/dploy/internal/utils"

type Task struct {
	Role    string      `json:"role"`
	rawTags interface{} `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func (task *Task) Tags() []string {
	return utils.InferSlice(task.rawTags)
}
