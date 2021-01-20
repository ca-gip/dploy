package ansible

import (
	"github.com/ca-gip/dploy/internal/utils"
)

type Play struct {
	Hosts   string      `json:"hosts" yaml:"hosts"`
	Roles   []Role      `yaml:"roles,omitempty"`
	RawTags interface{} `json:"tags,inline" yaml:"tags,inline"`
}

func (play *Play) Tags() []string {
	return utils.InferSlice(play.RawTags)
}
