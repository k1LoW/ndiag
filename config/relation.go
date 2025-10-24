package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elliotchance/orderedmap"
)

type RelationType struct {
	Name          string
	ComponentsKey string
	Attrs         Attrs
}

var RelationTypeDefault = &RelationType{
	Name:          "relation",
	ComponentsKey: "components",
	Attrs: Attrs{
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
	Attrs: Attrs{
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
	Components Components
	Labels     Labels
	Attrs      Attrs
}

func (r *Relation) ElementType() ElementType {
	return TypeRelation
}

func (r *Relation) FullName() string {
	return r.relationId
}

func (r *Relation) Id() string {
	return strings.ToLower(r.relationId)
}

func (r *Relation) DescFilename() string {
	return MakeMdFilename("_relation", r.Id())
}

type Relations []*Relation

func (relations Relations) FindByLabels(labels Labels) Relations {
	filtered := Relations{}
	for _, r := range relations {
		if len(r.Labels.Subtract(labels)) > 0 {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

type rawRelation struct {
	RelationId string `json:"id"`
	Desc       string
	Type       *RelationType
	Components []string
	Labels     []string `json:"-"`
	Attrs      Attrs
}

func (rel *rawRelation) Id() string {
	if rel.RelationId != "" {
		return strings.ToLower(rel.RelationId)
	}
	h := sha256.New()
	seed := []string{}
	for _, c := range rel.Components {
		seed = append(seed, queryTrim(c))
	}
	key := strings.ToLower(strings.Join(seed, "-"))
	if _, err := io.WriteString(h, string(key)); err != nil {
		return ""
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToLower(fmt.Sprintf("%s-%s", queryTrim(rel.Components[0]), s[:7]))
}

type rawRelations []*rawRelation

func (relations rawRelations) FindById(id string) (*rawRelation, error) {
	for _, r := range relations {
		if r.Id() == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("raw relation not found: %s", id)
}

func (dest rawRelations) Merge(src rawRelations) rawRelations {
	for _, sr := range src {
		r, err := dest.FindById(sr.Id())
		if err != nil {
			dest = append(dest, sr)
			continue
		}
		if sr.Desc != "" {
			r.Desc = sr.Desc
		}
		if sr.Type != nil {
			r.Type = sr.Type
		}
		r.Components = merge(r.Components, sr.Components)
		r.Labels = merge(r.Labels, sr.Labels)
		r.Attrs = r.Attrs.Merge(sr.Attrs)
	}
	return dest
}

func SplitRelations(relations Relations) []*Edge {
	var prev *Component
	edges := []*Edge{}
	for _, rel := range relations {
		prev = nil
		for _, r := range rel.Components {
			if prev != nil {
				edge := &Edge{
					Src:      prev,
					Dst:      r,
					Relation: rel,
					Attrs:    rel.Attrs,
				}
				prev.Edges = append(prev.Edges, edge)
				r.Edges = append(r.Edges, edge)
				edges = append(edges, edge)
			}
			prev = r
		}
	}
	return edges
}

func MergeEdges(edges []*Edge) []*Edge {
	eKeys0 := orderedmap.NewOrderedMap()
	merged0 := []*Edge{}
	for _, e := range edges {
		eKeys0.Set(fmt.Sprintf("%s->%s", e.Src.Id(), e.Dst.Id()), e)
	}
	for _, k := range eKeys0.Keys() {
		e, _ := eKeys0.Get(k)
		edge, ok := e.(*Edge)
		if !ok {
			continue
		}
		merged0 = append(merged0, edge)
	}

	eKeys1 := orderedmap.NewOrderedMap()
	merged1 := []*Edge{}
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
		edge, ok := e.(*Edge)
		if !ok {
			continue
		}
		merged1 = append(merged1, edge)
	}
	return merged1
}

func uniqueRawRelations(rels rawRelations) rawRelations {
	rKeys := orderedmap.NewOrderedMap()
	result := rawRelations{}
	for _, rel := range rels {
		key, _ := json.Marshal(rel)
		rKeys.Set(string(key), rel)
	}
	for _, k := range rKeys.Keys() {
		rel, _ := rKeys.Get(k)
		r, ok := rel.(*rawRelation)
		if !ok {
			continue
		}
		result = append(result, r)
	}
	return result
}
