package config

import (
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
	RealNodes     []*RealNode  `yaml:"-"`
	Labels        Labels
	nameRe        *regexp.Regexp
	rawComponents []string
	rawClusters   []string
}

type NodeMetadata struct {
	Icon   string   `yaml:"icon,omitempty"`
	Labels []string `yaml:"labels,omitempty"`
}

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

type RealNode struct {
	Node
	BelongTo *Node
}
