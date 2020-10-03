package config

import (
	"crypto/sha256"
	"fmt"
	"io"
	"regexp"
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

	for _, rel := range raw.Relations {
		components := []string{}
		tags := []string{}
		switch v := rel.(type) {
		case []interface{}:
			for _, r := range v {
				components = append(components, r.(string))
			}
			if len(components) < 2 {
				return fmt.Errorf("invalid relation format: %s", v)
			}
			id, err := genRelationId(components)
			if err != nil {
				return err
			}
			tags = []string{id}
			rrel := &rawRelation{
				Id:    id,
				Components: components,
				Tags:  tags,
			}
			d.rawRelations = append(d.rawRelations, rrel)
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
					return err
				}
			}
			ri, ok := v["components"]
			if !ok {
				return fmt.Errorf("invalid relation format: %s", v)
			}
			for _, r := range ri.([]interface{}) {
				components = append(components, r.(string))
			}
			if len(components) < 2 {
				return fmt.Errorf("invalid relation format: %s", v)
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
			rrel := &rawRelation{
				Id:    id,
				Components: components,
				Tags:  tags,
			}
			d.rawRelations = append(d.rawRelations, rrel)
		}
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

func genRelationId(components []string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, fmt.Sprintf("%s", components)); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("rel-%s", s[:12]), nil
}
