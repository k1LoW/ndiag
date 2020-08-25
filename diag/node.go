package diag

import (
	"regexp"
	"strings"
)

type Node struct {
	Name        string       `yaml:"name"`
	Desc        string       `yaml:"desc"`
	Components  []*Component `yaml:"components,omitempty"`
	Clusters    Clusters     `yaml:"clusters,omitempty"`
	RealNodes   []*RealNode
	nameRe      *regexp.Regexp
	rawClusters []string
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
