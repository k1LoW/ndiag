package config

import (
	"fmt"
	"strings"
)

type View struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc,omitempty"`
	Layers []string `yaml:"layers"`
	Labels []string `yaml:"labels,omitempty"`
}

func (v *View) ElementType() ElementType {
	return TypeView
}

func (v *View) FullName() string {
	switch {
	case v.Name != "":
		return v.Name
	case len(v.Layers) > 0 && len(v.Labels) > 0:
		return fmt.Sprintf("%s-%s", strings.Join(v.Layers, "-"), strings.Join(v.Labels, "-"))
	case len(v.Layers) > 0:
		return strings.Join(v.Layers, "-")
	case len(v.Labels) > 0:
		return strings.Join(v.Labels, "-")
	default:
		return ""
	}
}

func (v *View) Id() string {
	return strings.ToLower(v.FullName())
}

func (v *View) DescFilename() string {
	return MakeMdFilename("_view", v.Id())
}
