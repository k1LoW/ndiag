package config

import (
	"fmt"
	"strings"
)

type Component struct {
	Name     string
	Desc     string
	Cluster  *Cluster
	Node     *Node
	NEdges   []*NEdge
	Metadata ComponentMetadata
}

type ComponentMetadata struct {
	Icon     string `qs:"icon"`
	IconPath string `qs:"-"`
}

func (c *Component) FullName() string {
	switch {
	case c.Node != nil:
		// node components
		return fmt.Sprintf("%s:%s", c.Node.FullName(), c.Name)
	case c.Cluster != nil:
		// cluster components
		return fmt.Sprintf("%s:%s", c.Cluster.FullName(), c.Name)
	default:
		// global components
		return c.Name
	}
}

func (c *Component) Id() string {
	return strings.ToLower(c.FullName())
}

func (c *Component) OverrideMetadata(c2 *Component) error {
	if c.Id() != c2.Id() {
		return fmt.Errorf("can not merge: %s <-> %s", c.Id(), c2.Id())
	}
	if c2.Metadata.Icon != "" {
		c.Metadata.Icon = c2.Metadata.Icon
		c.Metadata.IconPath = c2.Metadata.IconPath
	}
	return nil
}
