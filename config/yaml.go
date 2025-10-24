package config

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/k1LoW/glyph"
	"github.com/k1LoW/tbls/dict"
)

func (d *Config) UnmarshalYAML(data []byte) error {
	raw := struct {
		Name          string             `yaml:"name"`
		Desc          string             `yaml:"desc,omitempty"`
		DocPath       string             `yaml:"docPath"`
		DescPath      string             `yaml:"descPath"`
		IconPath      string             `yaml:"iconPath,omitempty"`
		Graph         *Graph             `yaml:"graph,omitempty"`
		HideViews     bool               `yaml:"hideViews"`
		HideLayers    bool               `yaml:"hideLayers"`
		HideRealNodes bool               `yaml:"hideRealNodes"`
		HideLabels    bool               `yaml:"hideLabels"`
		Views         []*View            `yaml:"views"`
		Nodes         []*Node            `yaml:"nodes"`
		Networks      []interface{}      `yaml:"networks"`
		Relations     []interface{}      `yaml:"relations"`
		Dict          *dict.Dict         `yaml:"dict,omitempty"`
		BaseColor     string             `yaml:"baseColor,omitempty"`
		TextColor     string             `yaml:"textColor,omitempty"`
		CustomIcons   []*glyph.Blueprint `yaml:"customIcons,omitempty"`
	}{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Name = raw.Name
	d.Desc = raw.Desc
	d.DocPath = raw.DocPath
	d.DescPath = raw.DescPath
	d.IconPath = raw.IconPath
	if raw.Graph != nil {
		d.Graph = raw.Graph
	}
	d.HideViews = raw.HideViews
	d.HideLayers = raw.HideLayers
	d.HideRealNodes = raw.HideRealNodes
	d.HideLabels = raw.HideLabels
	d.Views = raw.Views
	d.Nodes = raw.Nodes
	if raw.Dict != nil {
		d.Dict = raw.Dict
	}
	d.BaseColor = raw.BaseColor
	d.TextColor = raw.TextColor
	d.CustomIcons = raw.CustomIcons

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
		Name        string       `yaml:"name"`
		Desc        string       `yaml:"desc"`
		Match       string       `yaml:"match,omitempty"`
		MatchRegexp string       `yaml:"matchRegexp,omitempty"`
		Components  []string     `yaml:"components,omitempty"`
		Clusters    []string     `yaml:"clusters,omitempty"`
		Metadata    NodeMetadata `yaml:"metadata,omitempty"`
	}{}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	if sepContains(raw.Name) {
		return fmt.Errorf("a node's name cannot contain unescaped '%s': %s ", Sep, raw.Name)
	}

	n.Name = raw.Name
	n.Match = raw.Match
	n.MatchRegexp = raw.MatchRegexp
	if n.Match == "" {
		n.Match = n.Name
	}
	if n.MatchRegexp == "" {
		n.nameRe = regexp.MustCompile(fmt.Sprintf("^%s$", strings.ReplaceAll(n.Match, "*", ".+")))
	} else {
		n.nameRe = regexp.MustCompile(n.MatchRegexp)
	}

	n.Desc = raw.Desc
	n.rawComponents = raw.Components
	n.rawClusters = raw.Clusters
	n.Metadata = raw.Metadata
	return nil
}

func (g *Graph) UnmarshalYAML(data []byte) error {
	raw := struct {
		Format        string        `yaml:"format,omitempty"`
		MapSliceAttrs yaml.MapSlice `yaml:"attrs,omitempty"`
	}{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return err
	}
	g.Format = raw.Format
	g.Attrs = Attrs{}
	for _, kv := range raw.MapSliceAttrs {
		key, ok := kv.Key.(string)
		if !ok {
			continue
		}
		value, ok := kv.Value.(string)
		if !ok {
			continue
		}
		g.Attrs = append(g.Attrs, &Attr{
			Key:   key,
			Value: value,
		})
	}
	return nil
}

func parseRelation(relType *RelationType, rel interface{}) (*rawRelation, error) {
	components := []string{}
	labels := []string{}
	switch v := rel.(type) {
	case []interface{}:
		// networks:
		//   - ["internet", "lb:nginx", "app:nginx", "app:rails"]
		for _, r := range v {
			str, ok := r.(string)
			if !ok {
				continue
			}
			components = append(components, str)
		}
		rel := &rawRelation{
			Type:       relType,
			Components: components,
			Attrs:      relType.Attrs,
		}
		rel.Labels = []string{}
		return rel, nil
	case map[string]interface{}:
		var (
			id string
		)
		idi, ok := v["id"]
		if ok {
			id, ok = idi.(string)
			if !ok {
				id = ""
			}
		} else {
			id = ""
		}
		ri, ok := v[relType.ComponentsKey]
		if ok {
			riSlice, ok := ri.([]interface{})
			if ok {
				for _, r := range riSlice {
					str, ok := r.(string)
					if !ok {
						continue
					}
					components = append(components, str)
				}
			}
		}
		typei, ok := v["type"]
		if ok {
			typeStr, ok := typei.(string)
			if ok {
				switch typeStr {
				case "network":
					relType = RelationTypeNetwork
				default:
					return nil, fmt.Errorf("invalid %s format: %s", relType.Name, v)
				}
			}
		}
		ti, ok := v["labels"]
		if ok {
			tiSlice, ok := ti.([]interface{})
			if ok {
				for _, t := range tiSlice {
					str, ok := t.(string)
					if !ok {
						continue
					}
					labels = append(labels, str)
				}
			}
		}
		attrs := Attrs{}
		attrsi, ok := v["attrs"]
		if ok {
			attrsMap, ok := attrsi.(map[string]interface{})
			if ok {
				for k, v := range attrsMap {
					vStr, ok := v.(string)
					if !ok {
						continue
					}
					attrs = append(attrs, &Attr{
						Key:   k,
						Value: vStr,
					})
				}
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
			RelationId: id,
			Type:       relType,
			Components: components,
			Labels:     labels,
			Attrs:      attrs,
		}, nil
	default:
		return nil, fmt.Errorf("invalid relation format: %s", v)
	}
}
