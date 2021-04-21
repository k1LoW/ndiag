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

type Views []*View

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

func (views Views) FindById(id string) (*View, error) {
	for _, v := range views {
		if v.Id() == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("view not found: %s", id)
}

func (dest Views) Merge(src Views) Views {
	for _, sv := range src {
		v, err := dest.FindById(sv.Id())
		if err != nil {
			dest = append(dest, sv)
			continue
		}
		if sv.Name != "" {
			v.Name = sv.Name
		}
		if sv.Desc != "" {
			v.Desc = sv.Desc
		}
		v.Layers = sv.Layers
		v.Labels = sv.Labels
	}
	return dest
}
