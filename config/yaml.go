package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

func (d *Config) UnmarshalYAML(data []byte) error {
	raw := struct {
		Name      string        `yaml:"name"`
		Desc      string        `yaml:"desc,omitempty"`
		DocPath   string        `yaml:"docPath"`
		Diagrams  []*Diagram    `yaml:"diagrams"`
		Nodes     []*Node       `yaml:"nodes"`
		Networks  []interface{} `yaml:"networks"`
		Relations []interface{} `yaml:"relations"`
	}{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Name = raw.Name
	d.Desc = raw.Desc
	d.DocPath = raw.DocPath
	d.Diagrams = raw.Diagrams
	d.Nodes = raw.Nodes

	for _, rel := range raw.Networks {
		rel, err := parseRelation(RelationTypeNetwork, rel)
		if err != nil {
			return err
		}
		d.rawRelations = append(d.rawRelations, rel)
	}

	for _, rel := range raw.Relations {
		rel, err := parseRelation(RelationTypeDefault, rel)
		if err != nil {
			return err
		}
		d.rawRelations = append(d.rawRelations, rel)
	}
	return nil
}

func (n *Node) UnmarshalYAML(data []byte) error {
	raw := struct {
		Name       string   `yaml:"name"`
		Desc       string   `yaml:"desc"`
		Components []string `yaml:"components,omitempty"`
		Clusters   []string `yaml:"clusters,omitempty"`
	}{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	if sepContains(raw.Name) {
		return fmt.Errorf("a node's name cannot contain unescaped '%s': %s ", Sep, raw.Name)
	}

	n.Name = raw.Name
	n.nameRe = regexp.MustCompile(fmt.Sprintf("^%s$", strings.Replace(n.Name, "*", ".+", -1)))
	n.Desc = raw.Desc
	n.rawComponents = raw.Components
	n.rawClusters = raw.Clusters
	return nil
}

func parseRelation(relType *RelationType, rel interface{}) (*rawRelation, error) {
	components := []string{}
	tags := []string{}
	switch v := rel.(type) {
	case []interface{}:
		for _, r := range v {
			components = append(components, r.(string))
		}
		if len(components) < 2 {
			return nil, fmt.Errorf("invalid %s format: %s", relType.Name, v)
		}
		rel := &rawRelation{
			Type:       relType,
			Components: components,
			Attrs:      relType.Attrs,
		}
		rel.Tags = []string{rel.Id()}
		return rel, nil
	case map[string]interface{}:
		var (
			id string
		)
		idi, ok := v["id"]
		if ok {
			id = idi.(string)
		} else {
			id = ""
		}
		ri, ok := v[relType.ComponentsKey]
		if !ok {
			return nil, fmt.Errorf("invalid %s format: %s", relType.Name, v)
		}
		for _, r := range ri.([]interface{}) {
			components = append(components, r.(string))
		}
		if len(components) < 2 {
			return nil, fmt.Errorf("invalid %s format: %s", relType.Name, v)
		}
		typei, ok := v["type"]
		if ok {
			switch typei.(string) {
			case "network":
				relType = RelationTypeNetwork
			default:
				return nil, fmt.Errorf("invalid %s format: %s", relType.Name, v)
			}
		}
		ti, ok := v["tags"]
		if ok {
			for _, t := range ti.([]interface{}) {
				tags = append(tags, t.(string))
			}
		}
		if len(tags) == 0 {
			tags = []string{id}
		}
		attrs := []*Attr{}
		attrsi, ok := v["attrs"]
		if ok {
			for k, v := range attrsi.(map[string]interface{}) {
				attrs = append(attrs, &Attr{
					Key:   k,
					Value: v.(string),
				})
			}
		}
		sort.Slice(attrs, func(i, j int) bool {
			if attrs[i].Key == attrs[j].Key {
				return attrs[i].Value < attrs[j].Value
			}
			return attrs[i].Key < attrs[j].Key
		})
		attrs = append(relType.Attrs, attrs...)

		return &rawRelation{
			relationId: id,
			Type:       relType,
			Components: components,
			Tags:       tags,
			Attrs:      attrs,
		}, nil
	default:
		return nil, fmt.Errorf("invalid relation format: %s", v)
	}
}
