package config

import (
	"fmt"
	"strings"

	"github.com/elliotchance/orderedmap"
)

type Network struct {
	NetworkId string
	Route     []*Component
	Tags      []string
}

func (n *Network) FullName() string {
	return fmt.Sprintf(n.NetworkId)
}

func (n *Network) Id() string {
	return strings.ToLower(n.FullName())
}

type rawNetwork struct {
	Id    string
	Route []string
	Tags  []string
}

type Tag struct {
	Name     string
	Desc     string
	Networks []*Network
}

func (t *Tag) FullName() string {
	return t.Name
}

func (t *Tag) Id() string {
	return strings.ToLower(t.FullName())
}

func SplitNetworks(networks []*Network) []*NEdge {
	var prev *Component
	edges := []*NEdge{}
	for _, nw := range networks {
		prev = nil
		for _, r := range nw.Route {
			if prev != nil {
				edge := &NEdge{
					Src:     prev,
					Dst:     r,
					Network: nw,
				}
				prev.NEdges = append(prev.NEdges, edge)
				r.NEdges = append(r.NEdges, edge)
				edges = append(edges, edge)
			}
			prev = r
		}
	}
	return edges
}

func MergeEdges(edges []*NEdge) []*NEdge {
	eKeys0 := orderedmap.NewOrderedMap()
	merged0 := []*NEdge{}
	for _, e := range edges {
		eKeys0.Set(fmt.Sprintf("%s->%s", e.Src.Id(), e.Dst.Id()), e)
	}
	for _, k := range eKeys0.Keys() {
		e, _ := eKeys0.Get(k)
		merged0 = append(merged0, e.(*NEdge))
	}

	eKeys1 := orderedmap.NewOrderedMap()
	merged1 := []*NEdge{}
	for _, e := range merged0 {
		var k string
		if e.Src.Id() < e.Dst.Id() {
			k = fmt.Sprintf("%s->%s", e.Src.Id(), e.Dst.Id())
		} else {
			k = fmt.Sprintf("%s->%s", e.Dst.Id(), e.Src.Id())
		}
		ce, _ := eKeys1.Get(k)
		if ce != nil {
			e.Attrs = append(e.Attrs, &Attr{
				Key:   "dir",
				Value: "both",
			})
		}
		eKeys1.Set(k, e)
	}
	for _, k := range eKeys1.Keys() {
		e, _ := eKeys1.Get(k)
		merged1 = append(merged1, e.(*NEdge))
	}
	return merged1
}
