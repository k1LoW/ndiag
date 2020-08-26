package diag

import (
	"fmt"
	"strings"
)

type Cluster struct {
	Layer      string
	Name       string
	Parent     *Cluster
	Children   []*Cluster
	Nodes      []*Node
	Components []*Component
}

func (c *Cluster) FullName() string {
	return fmt.Sprintf("%s:%s", c.Layer, c.Name)
}

func (c *Cluster) Id() string {
	return strings.ToLower(c.FullName())
}

type Clusters []*Cluster

func (cs Clusters) Find(layer, name string) *Cluster {
	for _, c := range cs {
		if c.Layer == layer && c.Name == name {
			return c
		}
	}
	return nil
}

func (cs Clusters) FindByLayer(layer string) Clusters {
	result := Clusters{}
	for _, c := range cs {
		if c.Layer == layer {
			result = append(result, c)
		}
	}
	return result
}
