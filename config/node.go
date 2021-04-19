package config

import (
	"fmt"
	"regexp"
	"strings"
)

type Node struct {
	Name          string       `yaml:"name"`
	Desc          string       `yaml:"desc,omitempty"`
	Match         string       `yaml:"match,omitempty"`
	MatchRegexp   string       `yaml:"matchRegexp,omitempty"`
	Components    []*Component `yaml:"components,omitempty"`
	Clusters      Clusters     `yaml:"clusters,omitempty"`
	Metadata      NodeMetadata `yaml:"metadata,omitempty"`
	RealNodes     RealNodes    `yaml:"-"`
	Labels        Labels
	nameRe        *regexp.Regexp
	rawComponents []string
	rawClusters   []string
}

type NodeMetadata struct {
	Icon   string   `yaml:"icon,omitempty"`
	Labels []string `yaml:"labels,omitempty"`
}

type Nodes []*Node

type RealNode struct {
	Node
	BelongTo *Node
}

type RealNodes []*RealNode

func (n *Node) ElementType() ElementType {
	return TypeNode
}

func (n *Node) FullName() string {
	return n.Name
}

func (n *Node) Id() string {
	return strings.ToLower(n.FullName())
}

func (n *Node) DescFilename() string {
	return MakeMdFilename("_node", n.Id())
}

func (nodes Nodes) FindById(id string) (*Node, error) {
	for _, v := range nodes {
		if v.Id() == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("node not found: %s", id)
}

func (nodes RealNodes) FindById(id string) (*RealNode, error) {
	for _, v := range nodes {
		if v.Id() == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("real node not found: %s", id)
}
