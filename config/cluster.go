package config

import (
	"fmt"
	"image/color"
	"strings"
)

type Cluster struct {
	Layer      *Layer
	Name       string
	Desc       string
	Parent     *Cluster
	Children   Clusters
	Nodes      Nodes
	Components Components
	Metadata   ClusterMetadata
}

type ClusterMetadata struct {
	Icon      string      `qs:"icon"`
	Color     color.Color `qs:"-"`
	FillColor color.Color `qs:"-"`
	TextColor color.Color `qs:"-"`
}

func (c *Cluster) ElementType() ElementType {
	return TypeCluster
}

func (c *Cluster) FullName() string {
	return fmt.Sprintf("%s:%s", c.Layer.FullName(), c.Name)
}

func (c *Cluster) Id() string {
	return strings.ToLower(c.FullName())
}

func (c *Cluster) DescFilename() string {
	return MakeMdFilename("_cluster", c.Id())
}

func (c *Cluster) OverrideMetadata(c2 *Cluster) error {
	if c.Id() != c2.Id() {
		return fmt.Errorf("can not merge: %s <-> %s", c.Id(), c2.Id())
	}
	if c2.Metadata.Icon != "" {
		c.Metadata.Icon = c2.Metadata.Icon
	}
	if c2.Metadata.Color != nil {
		c.Metadata.Color = c2.Metadata.Color
		c.Metadata.FillColor = c2.Metadata.FillColor
		c.Metadata.TextColor = c2.Metadata.TextColor
	}
	return nil
}

type Clusters []*Cluster

func (cs Clusters) Find(layer, name string) *Cluster {
	for _, c := range cs {
		if strings.EqualFold(c.Layer.Id(), layer) && c.Name == name {
			return c
		}
	}
	return nil
}

func (cs Clusters) FindById(id string) (*Cluster, error) {
	for _, c := range cs {
		if c.Layer.Id() == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("cluster not found: %s", id)
}

func (cs Clusters) FindByLayer(layer string) Clusters {
	result := Clusters{}
	for _, c := range cs {
		if strings.EqualFold(c.Layer.Id(), layer) {
			result = append(result, c)
		}
	}
	return result
}

func (cs Clusters) Contains(t *Cluster) bool {
	for _, c := range cs {
		if c.Id() == t.Id() {
			return true
		}
		if c.Children.Contains(t) {
			return true
		}
	}
	return false
}
