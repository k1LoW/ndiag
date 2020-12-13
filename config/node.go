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
	nameRe        *regexp.Regexp
	rawComponents []string
	rawClusters   []string
}

type NodeMetadata struct {
	Icon     string `yaml:"icon,omitempty"`
	IconPath string `yaml:"-"`
}

func (n *Node) FullName() string {
	return n.Name
}

func (n *Node) Id() string {
	return strings.ToLower(n.FullName())
}

type RealNode struct {
	Node
	BelongTo *Node
}
