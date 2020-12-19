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
	Children   []*Cluster
	Nodes      []*Node
	Components []*Component
	Metadata   ClusterMetadata
}

type ClusterMetadata struct {
	Color     color.Color
	FillColor color.Color
	TextColor color.Color
}

func (c *Cluster) FullName() string {
	return fmt.Sprintf("%s:%s", c.Layer.Name, c.Name)
}

func (c *Cluster) Id() string {
	return strings.ToLower(c.FullName())
}

type Clusters []*Cluster

func (cs Clusters) Find(layer, name string) *Cluster {
	for _, c := range cs {
		if c.Layer.Name == layer && c.Name == name {
			return c
		}
	}
	return nil
}

func (cs Clusters) FindByLayer(layer string) Clusters {
	result := Clusters{}
	for _, c := range cs {
		if c.Layer.Name == layer {
			result = append(result, c)
		}
	}
	return result
}
