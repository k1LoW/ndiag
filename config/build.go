package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/orderedmap"
	"github.com/pasztorpisti/qs"
)

func (cfg *Config) buildDefault() error {
	if cfg.DocPath == "" {
		cfg.DocPath = DefaultDocPath
	}
	if !filepath.IsAbs(cfg.DocPath) {
		docPath, err := filepath.Abs(filepath.Join(cfg.basePath, cfg.DocPath))
		if err != nil {
			return err
		}
		cfg.DocPath = docPath
	}
	if cfg.DescPath == "" {
		cfg.DescPath = DefaultDescPath
	}
	if !filepath.IsAbs(cfg.DescPath) {
		descPath, err := filepath.Abs(filepath.Join(cfg.basePath, cfg.DescPath))
		if err != nil {
			return err
		}
		cfg.DescPath = descPath
	}
	if cfg.IconPath == "" {
		cfg.IconPath = DefaultIconPath
	}
	if !filepath.IsAbs(cfg.IconPath) {
		iconPath, err := filepath.Abs(filepath.Join(cfg.basePath, cfg.IconPath))
		if err != nil {
			return err
		}
		cfg.IconPath = iconPath
	}
	if cfg.BaseColor == "" {
		cfg.BaseColor = DefaultBaseColor
	}
	if cfg.TextColor == "" {
		cfg.TextColor = DefaultTextColor
	}
	return nil
}

func (cfg *Config) buildNodes() error {
	for _, n := range cfg.Nodes {
		if n.Metadata.Icon != "" {
			_, err := cfg.IconMap().Get(n.Metadata.Icon)
			if err != nil {
				return fmt.Errorf("not found icon: %s", n.Metadata.Icon)
			}
		}
		for _, s := range n.Metadata.Labels {
			n.Labels = append(n.Labels, cfg.FindOrCreateLabel(s))
		}
	}
	return nil
}

func (cfg *Config) buildComponents() error {
	gc := orderedmap.NewOrderedMap()
	nc := orderedmap.NewOrderedMap()
	cc := orderedmap.NewOrderedMap()
	for _, rel := range cfg.rawRelations {
		for _, r := range rel.Components {
			switch sepCount(r) {
			case 2: // cluster components
				cc.Set(r, struct{}{})
			case 1: // node components
				nc.Set(r, struct{}{})
			case 0: // global components
				gc.Set(r, struct{}{})
			}
		}
	}

	// global components
	for _, c := range gc.Keys() {
		com, err := cfg.parseComponent(c.(string))
		if err != nil {
			return err
		}
		current, err := cfg.FindComponent(com.Id())
		if err != nil {
			// create global component from relations
			cfg.globalComponents = append(cfg.globalComponents, com)
		} else {
			if err := current.OverrideMetadata(com); err != nil {
				return err
			}
		}
	}

	// node components
	for _, n := range cfg.Nodes {
		for _, c := range n.rawComponents {
			com, err := cfg.parseComponent(c)
			if err != nil {
				return err
			}
			com.Node = n
			n.Components = append(n.Components, com)
		}
		cfg.nodeComponents = append(cfg.nodeComponents, n.Components...)
	}

	for _, c := range nc.Keys() {
		splitted := sepSplit(c.(string))
		nodeName := splitted[0]
		comName := splitted[1]
		n, err := cfg.FindNode(nodeName)
		if err != nil {
			return fmt.Errorf("node '%s' not found: %s", nodeName, c)
		}
		com, err := cfg.parseComponent(comName)
		if err != nil {
			return err
		}
		com.Node = n
		current, err := cfg.FindComponent(com.Id())
		if err != nil {
			// create node component from relations
			n.Components = append(n.Components, com)
			cfg.nodeComponents = append(cfg.nodeComponents, com)
		} else {
			if err := current.OverrideMetadata(com); err != nil {
				return err
			}
		}
	}

	// cluster components
	for _, c := range cc.Keys() {
		splitted := sepSplit(c.(string))
		clName := fmt.Sprintf("%s:%s", splitted[0], splitted[1])
		comName := splitted[2]
		belongTo := false
		for _, cl := range cfg.Clusters() {
			if strings.EqualFold(cl.FullName(), queryTrim(clName)) {
				com, err := cfg.parseComponent(comName)
				if err != nil {
					return err
				}
				com.Cluster = cl
				current, err := cfg.FindComponent(com.Id())
				if err != nil {
					// create cluster component from relations
					cl.Components = append(cl.Components, com)
					cfg.clusterComponents = append(cfg.clusterComponents, com)
				} else {
					if err := current.OverrideMetadata(com); err != nil {
						return err
					}
				}
				belongTo = true
				break
			}
		}
		if !belongTo {
			return fmt.Errorf("cluster '%s' not found: %s", clName, c)
		}
	}

	// find or create labels
	for _, c := range cfg.Components() {
		for _, l := range c.Metadata.Labels {
			c.Labels = append(c.Labels, cfg.FindOrCreateLabel(l))
		}
	}

	return nil
}

func (cfg *Config) buildClusters() error {
	for _, n := range cfg.Nodes {
		for _, c := range n.rawClusters {
			cluster, err := cfg.parseAndCollectCluster(c)
			if err != nil {
				return err
			}
			cluster.Nodes = append(cluster.Nodes, n)
			n.Clusters = append(n.Clusters, cluster)
		}
	}
	for _, rel := range cfg.rawRelations {
		for _, r := range rel.Components {
			if sepCount(r) != 2 {
				continue
			}
			splited := sepSplit(r)
			clusterId := sepJoin(splited[:2])
			if _, err := cfg.parseAndCollectCluster(clusterId); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfg *Config) parseAndCollectCluster(clusterId string) (*Cluster, error) {
	if !sepContains(clusterId) {
		return nil, fmt.Errorf("invalid cluster id: %s", clusterId)
	}
	if sepCount(clusterId) != 1 {
		return nil, fmt.Errorf("invalid cluster id: %s", clusterId)
	}
	splitted := sepSplit(clusterId)
	layerStr := splitted[0]
	name := splitted[1]

	var m ClusterMetadata
	if queryContains(name) {
		splited := querySplit(name)
		name = splited[0]
		if err := qs.Unmarshal(&m, splited[1]); err != nil {
			return nil, err
		}
		if m.Icon != "" {
			_, err := cfg.IconMap().Get(m.Icon)
			if err != nil {
				return nil, fmt.Errorf("not found icon: %s", m.Icon)
			}
		}
	}

	layer, err := cfg.FindLayer(layerStr)
	if err != nil {
		layer = &Layer{Name: layerStr}
		cfg.layers = append(cfg.layers, layer)
	}
	newC := &Cluster{
		Layer:    layer,
		Name:     name,
		Metadata: m,
	}

	current := cfg.clusters.Find(layerStr, name)
	if current != nil {
		if err := current.OverrideMetadata(newC); err != nil {
			return nil, err
		}
		return current, nil
	}

	cfg.clusters = append(cfg.clusters, newC)

	return newC, nil
}

func (cfg *Config) buildRelations() error {
	for _, rel := range cfg.rawRelations {
		labels := []*Label{}
		for _, s := range rel.Labels {
			labels = append(labels, cfg.FindOrCreateLabel(s))
		}
		nrel := &Relation{
			relationId: rel.Id(),
			Type:       rel.Type,
			Labels:     labels,
			Attrs:      rel.Attrs,
		}
		for _, r := range rel.Components {
			c, err := cfg.FindComponent(r)
			if err != nil {
				return err
			}
			nrel.Components = append(nrel.Components, c)
		}
		cfg.Relations = append(cfg.Relations, nrel)
	}
	cfg.edges = SplitRelations(cfg.Relations)

	return nil
}

func (cfg *Config) buildDescriptions() error {
	if cfg.DescPath == "" {
		return nil
	}
	err := os.MkdirAll(cfg.DescPath, 0755) // #nosec
	if err != nil {
		return err
	}

	// top
	if cfg.Desc == "" {
		desc, err := cfg.readDescFile(MakeMdFilename("_index", ""))
		if err != nil {
			return err
		}
		cfg.Desc = desc
	}

	// views
	if !cfg.HideViews {
		for _, v := range cfg.Views {
			if v.Desc != "" {
				continue
			}
			path := filepath.Join(cfg.DescPath, MakeMdFilename("_view", v.Id()))
			oldPath := filepath.Join(cfg.DescPath, MakeMdFilename("_diagram", v.Id()))
			if _, err := os.Stat(oldPath); err == nil {
				if _, err := os.Stat(path); err == nil {
					return fmt.Errorf("old description file exists: %s", oldPath)
				}
				if err := os.Rename(oldPath, path); err != nil {
					return err
				}
			}
			desc, err := cfg.readDescFile(MakeMdFilename("_view", v.Id()))
			if err != nil {
				return err
			}
			v.Desc = desc
		}
	}

	// clusters
	for _, c := range cfg.Clusters() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_cluster", c.Id()))
		if err != nil {
			return err
		}
		c.Desc = desc
	}

	// layers
	for _, l := range cfg.Layers() {
		if l.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_layer", l.Id()))
		if err != nil {
			return err
		}
		l.Desc = desc
	}

	// nodes
	for _, n := range cfg.Nodes {
		if n.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_node", n.Id()))
		if err != nil {
			return err
		}
		n.Desc = desc
	}

	// components
	for _, c := range cfg.GlobalComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_component", c.Id()))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.ClusterComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_component", c.Id()))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.NodeComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_component", c.Id()))
		if err != nil {
			return err
		}
		c.Desc = desc
	}

	// labels
	for _, l := range cfg.labels {
		if l.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_label", l.Id()))
		if err != nil {
			return err
		}
		l.Desc = desc
	}

	// relations
	for _, r := range cfg.Relations {
		if r.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MakeMdFilename("_relation", r.Id()))
		if err != nil {
			return err
		}
		r.Desc = desc
	}

	return nil
}

func (cfg *Config) readDescFile(f string) (string, error) {
	descPath := filepath.Join(cfg.DescPath, f)
	file, err := os.OpenFile(descPath, os.O_RDONLY|os.O_CREATE, 0644) // #nosec
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return string(b), err
}

func (cfg *Config) buildColors() error {
	cfg.colorSets = defaultColorSets(cfg.BaseColor, cfg.TextColor)

	// layers
	for i, l := range cfg.Layers() {
		if l.Metadata.Color == nil {
			l.Metadata.Color = cfg.colorSets.Get(i).Color
		}
		if l.Metadata.FillColor == nil {
			l.Metadata.FillColor = cfg.colorSets.Get(i).FillColor
		}
		if l.Metadata.TextColor == nil {
			l.Metadata.TextColor = cfg.colorSets.Get(i).TextColor
		}
	}

	// clusters
	for _, c := range cfg.Clusters() {
		if c.Metadata.Color == nil {
			c.Metadata.Color = c.Layer.Metadata.Color
		}
		if c.Metadata.FillColor == nil {
			c.Metadata.FillColor = c.Layer.Metadata.FillColor
		}
		if c.Metadata.TextColor == nil {
			c.Metadata.TextColor = c.Layer.Metadata.TextColor
		}
	}
	return nil
}
