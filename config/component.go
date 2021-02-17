package config

import (
	"fmt"
	"sort"
	"strings"
)

type Component struct {
	Name     string
	Desc     string
	Cluster  *Cluster
	Node     *Node
	Edges    []*Edge
	Labels   Labels
	Metadata ComponentMetadata
}

type ComponentMetadata struct {
	Icon   string   `qs:"icon"`
	Labels []string `qs:"label"`
}

func (c *Component) ElementType() ElementType {
	return TypeComponent
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

func (c *Component) DescFilename() string {
	return MakeMdFilename("_component", c.Id())
}

func (c *Component) OverrideMetadata(c2 *Component) error {
	if c.Id() != c2.Id() {
		return fmt.Errorf("can not merge: %s <-> %s", c.Id(), c2.Id())
	}
	if c2.Metadata.Icon != "" {
		c.Metadata.Icon = c2.Metadata.Icon
	}
	if len(c2.Metadata.Labels) > 0 {
		c.Metadata.Labels = uniqueAndSort(append(c.Metadata.Labels, c2.Metadata.Labels...))
	}
	return nil
}

type Components []*Component

func (components Components) FindById(id string) (*Component, error) {
	for _, c := range components {
		if c.Id() == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("component not found: %s", id)
}

func uniqueAndSort(in []string) []string {
	m := map[string]struct{}{}
	for _, s := range in {
		m[s] = struct{}{}
	}
	u := []string{}
	for s := range m {
		u = append(u, s)
	}
	sort.Slice(u, func(i, j int) bool {
		return u[i] < u[j]
	})
	return u
}
