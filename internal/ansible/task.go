package ansible

import "github.com/ca-gip/dploy/internal/utils"

type Task struct {
	Name string `yaml:"name,omitempty"`
	Tags utils.Set
}
