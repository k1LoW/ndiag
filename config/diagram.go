package config

type Diagram struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc,omitempty"`
	Layers []string `yaml:"layers"`
}

// NewDiagram
func NewDiagram(name, desc string, layers []string) *Diagram {
	return &Diagram{
		Name:   name,
		Desc:   desc,
		Layers: layers,
	}
}
