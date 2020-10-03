package config

import (
	"crypto/sha256"
	"fmt"
	"io"
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
		rel, err := parseRelation("network", "route", rel)
		if err != nil {
			return err
		}
		rel.Attrs = append(defaultNetworkAttrs, rel.Attrs...)
		d.rawRelations = append(d.rawRelations, rel)
	}

	for _, rel := range raw.Relations {
		rel, err := parseRelation("relation", "components", rel)
		if err != nil {
			return err
		}
		rel.Attrs = append(defaultRelationAttrs, rel.Attrs...)
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
	n.Components = []*Component{}
	for _, c := range raw.Components {
		n.Components = append(n.Components, &Component{
			Name: c,
			Node: n,
		})
	}
	n.rawClusters = raw.Clusters
	return nil
}

func parseRelation(relType, componentKey string, rel interface{}) (*rawRelation, error) {
	components := []string{}
	tags := []string{}
	switch v := rel.(type) {
	case []interface{}:
		for _, r := range v {
			components = append(components, r.(string))
		}
		if len(components) < 2 {
			return nil, fmt.Errorf("invalid %s format: %s", relType, v)
		}
		id, err := genRelationId(components)
		if err != nil {
			return nil, err
		}
		tags = []string{id}
		return &rawRelation{
			Id:         id,
			Components: components,
			Tags:       tags,
		}, nil
	case map[string]interface{}:
		var (
			id  string
			err error
		)
		idi, ok := v["id"]
		if ok {
			id = idi.(string)
		} else {
			id, err = genRelationId(components)
			if err != nil {
				return nil, err
			}
		}
		ri, ok := v[componentKey]
		if !ok {
			return nil, fmt.Errorf("invalid %s format: %s", relType, v)
		}
		for _, r := range ri.([]interface{}) {
			components = append(components, r.(string))
		}
		if len(components) < 2 {
			return nil, fmt.Errorf("invalid %s format: %s", relType, v)
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

		return &rawRelation{
			Id:         id,
			Components: components,
			Tags:       tags,
			Attrs:      attrs,
		}, nil
	default:
		return nil, fmt.Errorf("invalid relation format: %s", v)
	}
}

func genRelationId(components []string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, fmt.Sprintf("%s", components)); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("rel-%s", s[:12]), nil
}
