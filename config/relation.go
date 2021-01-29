package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/elliotchance/orderedmap"
)

type RelationType struct {
	Name          string
	ComponentsKey string
	Attrs         []*Attr
}

var RelationTypeDefault = &RelationType{
	Name:          "relation",
	ComponentsKey: "components",
	Attrs: []*Attr{
		&Attr{
			Key:   "color",
			Value: "#4B75B9",
		},
		&Attr{
			Key:   "arrowhead",
			Value: "dot",
		},
		&Attr{
			Key:   "arrowhead",
			Value: "dot",
		},
		&Attr{
			Key:   "style",
			Value: "bold,dashed",
		},
	},
}

var RelationTypeNetwork = &RelationType{
	Name:          "network",
	ComponentsKey: "route",
	Attrs: []*Attr{
		&Attr{
			Key:   "color",
			Value: "#33333399",
		},
		&Attr{
			Key:   "arrowhead",
			Value: "normal",
		},
		&Attr{
			Key:   "arrowhead",
			Value: "normal",
		},
		&Attr{
			Key:   "style",
			Value: "bold",
		},
	},
}

type Relation struct {
	relationId string
	Desc       string
	Type       *RelationType
	Components []*Component
	Labels     []*Label
	Attrs      []*Attr
}

func (n *Relation) FullName() string {
	return n.relationId
}

func (n *Relation) Id() string {
	return strings.ToLower(n.relationId)
}

type rawRelation struct {
	RelationId string `json:"id"`
	Desc       string
	Type       *RelationType
	Components []string
	Labels     []string `json:"-"`
	Attrs      []*Attr
}

func (rel *rawRelation) Id() string {
	if rel.RelationId != "" {
		return strings.ToLower(rel.RelationId)
	}
	h := sha256.New()
	seed := []string{}
	seed = append(seed, rel.Type.Name)
	for _, c := range rel.Components {
		seed = append(seed, queryTrim(c))
	}
	sort.Slice(rel.Labels, func(i, j int) bool {
		return rel.Labels[i] < rel.Labels[j]
	})
	for _, l := range rel.Labels {
		seed = append(seed, l)
	}
	key := strings.Join(seed, "-")
	if _, err := io.WriteString(h, string(key)); err != nil {
		return ""
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("%s-%s", queryTrim(rel.Components[0]), s[:7])
}

func SplitRelations(relations []*Relation) []*NEdge {
	var prev *Component
	edges := []*NEdge{}
	for _, rel := range relations {
		prev = nil
		for _, r := range rel.Components {
			if prev != nil {
				edge := &NEdge{
					Src:      prev,
					Dst:      r,
					Relation: rel,
					Attrs:    rel.Attrs,
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

func uniqueRawRelations(rels []*rawRelation) []*rawRelation {
	rKeys := orderedmap.NewOrderedMap()
	result := []*rawRelation{}
	for _, rel := range rels {
		key, _ := json.Marshal(rel)
		rKeys.Set(string(key), rel)
	}
	for _, k := range rKeys.Keys() {
		rel, _ := rKeys.Get(k)
		result = append(result, rel.(*rawRelation))
	}
	return result
}
