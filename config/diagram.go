package config

import (
	"fmt"
	"strings"
)

type Diagram struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc,omitempty"`
	Layers []string `yaml:"layers"`
	Tags   []string `yaml:"tags,omitempty"`
}

func (d *Diagram) FullName() string {
	switch {
	case d.Name != "":
		return d.Name
	case len(d.Layers) > 0 && len(d.Tags) > 0:
		return fmt.Sprintf("%s-%s", strings.Join(d.Layers, "-"), strings.Join(d.Tags, "-"))
	case len(d.Layers) > 0:
		return strings.Join(d.Layers, "-")
	case len(d.Tags) > 0:
		return strings.Join(d.Tags, "-")
	default:
		return ""
	}
}

func (d *Diagram) Id() string {
	return strings.ToLower(d.FullName())
}
