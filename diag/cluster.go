package diag

import (
	"fmt"
	"strings"
)

type Cluster struct {
	Key        string
	Name       string
	Parent     *Cluster
	Children   []*Cluster
	Nodes      []*Node
	Components []*Component
}

func (c *Cluster) FullName() string {
	return fmt.Sprintf("%s:%s", c.Key, c.Name)
}

func (c *Cluster) Id() string {
	return strings.ToLower(c.FullName())
}

type Clusters []*Cluster

func (cs Clusters) Find(key, name string) *Cluster {
	for _, c := range cs {
		if c.Key == key && c.Name == name {
			return c
		}
	}
	return nil
}

func (cs Clusters) FindByKey(key string) Clusters {
	result := Clusters{}
	for _, c := range cs {
		if c.Key == key {
			result = append(result, c)
		}
	}
	return result
}
