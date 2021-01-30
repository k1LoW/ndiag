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

func (d *View) FullName() string {
	switch {
	case d.Name != "":
		return d.Name
	case len(d.Layers) > 0 && len(d.Labels) > 0:
		return fmt.Sprintf("%s-%s", strings.Join(d.Layers, "-"), strings.Join(d.Labels, "-"))
	case len(d.Layers) > 0:
		return strings.Join(d.Layers, "-")
	case len(d.Labels) > 0:
		return strings.Join(d.Labels, "-")
	default:
		return ""
	}
}

func (d *View) Id() string {
	return strings.ToLower(d.FullName())
}
