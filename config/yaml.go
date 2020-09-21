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
		route := []string{}
		switch v := nw.(type) {
		case []interface{}:
			for _, r := range v {
				route = append(route, r.(string))
			}
			if len(route) < 2 {
				return fmt.Errorf("invalid network format: %s", v)
			}
			id, err := genNetworkId(route)
			if err != nil {
				return err
			}
			rnw := &rawNetwork{
				Id:    id,
				Route: route,
			}
			d.rawNetworks = append(d.rawNetworks, rnw)
		case map[string]interface{}:
			id, ok := v["id"]
			if !ok {
				return fmt.Errorf("invalid network format: %s", v)
			}
			ri, ok := v["route"]
			if !ok {
				return fmt.Errorf("invalid network format: %s", v)
			}
			for _, r := range ri.([]interface{}) {
				route = append(route, r.(string))
			}
			if len(route) < 2 {
				return fmt.Errorf("invalid network format: %s", v)
			}
			rnw := &rawNetwork{
				Id:    id.(string),
				Route: route,
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

func genNetworkId(route []string) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, fmt.Sprintf("%s", route)); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return fmt.Sprintf("nw-%s", s[:12]), nil
}
