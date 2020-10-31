package config

type Diagram struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc,omitempty"`
	Layers []string `yaml:"layers"`
	Tags   []string `yaml:"tags,omitempty"`
}
