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
	Components    Components   `yaml:"components,omitempty"`
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

func (n *Node) OverrideMetadata(n2 *Node) error {
	if n.Id() != n2.Id() {
		return fmt.Errorf("can not merge: %s <-> %s", n.Id(), n2.Id())
	}
	if n2.Metadata.Icon != "" {
		n.Metadata.Icon = n2.Metadata.Icon
	}
	if len(n2.Metadata.Labels) > 0 {
		n.Metadata.Labels = uniqueAndSort(append(n.Metadata.Labels, n2.Metadata.Labels...))
	}
	return nil
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

func (dest Nodes) Merge(src Nodes) (Nodes, error) {
	for _, sn := range src {
		if len(sn.Components) > 0 || len(sn.Clusters) > 0 || len(sn.RealNodes) > 0 || len(sn.Labels) > 0 {
			return nil, fmt.Errorf("it should be before the Config.Build phase: %s", sn.Id())
		}
		n, err := dest.FindById(sn.Id())
		if err != nil {
			dest = append(dest, sn)
			continue
		}
		if len(n.Components) > 0 || len(n.Clusters) > 0 || len(n.RealNodes) > 0 || len(n.Labels) > 0 {
			return nil, fmt.Errorf("it should be before the Config.Build phase: %s", n.Id())
		}
		if sn.Desc != "" {
			n.Desc = sn.Desc
		}
		if sn.Match != "" {
			n.Match = sn.Match
		}
		if sn.MatchRegexp != "" {
			n.MatchRegexp = sn.MatchRegexp
			n.nameRe = sn.nameRe
		}
		if err := n.OverrideMetadata(sn); err != nil {
			return nil, err
		}
		n.rawComponents = uniqueAndSort(append(n.rawComponents, sn.rawComponents...))
		n.rawClusters = uniqueAndSort(append(n.rawClusters, sn.rawClusters...))
	}
	return dest, nil
}
