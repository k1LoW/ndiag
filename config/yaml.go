package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

func (d *Config) UnmarshalYAML(data []byte) error {
	raw := struct {
		Name     string        `yaml:"name"`
		Desc     string        `yaml:"desc,omitempty"`
		DocPath  string        `yaml:"docPath"`
		Diagrams []*Diagram    `yaml:"diagrams"`
		Nodes    []*Node       `yaml:"nodes"`
		Networks []interface{} `yaml:"networks"`
	}{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Name = raw.Name
	d.Desc = raw.Desc
	d.DocPath = raw.DocPath
	d.Diagrams = raw.Diagrams
	d.Nodes = raw.Nodes

	for _, nw := range raw.Networks {
		switch v := nw.(type) {
		case []interface{}:
			if len(v) != 2 {
				return fmt.Errorf("invalid network format: %s", v)
			}
			rnw := &rawNetwork{
				Src: v[0].(string),
				Dst: v[1].(string),
			}
			if len(v) == 3 {
				rnw.Desc = v[2].(string)
			}
			d.rawNetworks = append(d.rawNetworks, rnw)
		case map[string]interface{}:
			src, ok := v["src"]
			if !ok {
				return fmt.Errorf("invalid network format: %s", v)
			}
			dst, ok := v["dst"]
			if !ok {
				return fmt.Errorf("invalid network format: %s", v)
			}
			rnw := &rawNetwork{
				Src: src.(string),
				Dst: dst.(string),
			}
			if desc, ok := v["desc"]; ok {
				rnw.Desc = desc.(string)
			}
			d.rawNetworks = append(d.rawNetworks, rnw)
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
