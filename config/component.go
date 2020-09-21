package config

import (
	"fmt"
	"strings"
)

type Component struct {
	Name    string
	Desc    string
	Cluster *Cluster
	Node    *Node
	NEdges  []*NEdge
}

func (c *Component) FullName() string {
	if c.Node == nil {
		if c.Cluster == nil {
			// global components
			return c.Name
		}
		// cluster components
		return fmt.Sprintf("%s:%s", c.Cluster.FullName(), c.Name)
	}
	// node components
	return fmt.Sprintf("%s:%s", c.Node.FullName(), c.Name)
}

func (c *Component) Id() string {
	return strings.ToLower(c.FullName())
}
