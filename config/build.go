package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/orderedmap"
)

func (cfg *Config) buildNodes() error {
	for _, n := range cfg.Nodes {
		if n.Metadata.Icon != "" {
			n.Metadata.IconPath = filepath.Join(cfg.TempIconDir(), fmt.Sprintf("%s.%s", n.Metadata.Icon, cfg.Format()))
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
			if strings.EqualFold(cl.FullName(), clName) {
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
	return nil
}

func (cfg *Config) buildClusters() error {
	for _, n := range cfg.Nodes {
		for _, c := range n.rawClusters {
			cluster, err := cfg.parseClusterLabel(c)
			if err != nil {
				return err
			}
			cluster.Nodes = append(cluster.Nodes, n)
			n.Clusters = append(n.Clusters, cluster)
		}
	}
	return nil
}

func (cfg *Config) parseClusterLabel(label string) (*Cluster, error) {
	if !strings.Contains(label, Sep) {
		return nil, fmt.Errorf("invalid cluster id: %s", label)
	}
	splitted := sepSplit(label)
	if len(splitted) != 2 {
		return nil, fmt.Errorf("invalid cluster id: %s", label)
	}
	layerStr := splitted[0]
	name := splitted[1]
	current := cfg.clusters.Find(layerStr, name)
	if current != nil {
		return current, nil
	}
	layer, err := cfg.FindLayer(layerStr)
	if err != nil {
		layer = &Layer{Name: layerStr}
		cfg.layers = append(cfg.layers, layer)
	}
	newC := &Cluster{
		Layer: layer,
		Name:  name,
	}
	cfg.clusters = append(cfg.clusters, newC)

	return newC, nil
}

func (cfg *Config) buildRelations() error {
	relLabels := orderedmap.NewOrderedMap()
	for _, rel := range cfg.rawRelations {
		nrel := &Relation{
			relationId: rel.Id(),
			Type:       rel.Type,
			Labels:     rel.Labels,
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

		// labels
		for _, t := range rel.Labels {
			if t == "" {
				continue
			}
			var nt *Label
			nti, ok := relLabels.Get(t)
			if ok {
				nt = nti.(*Label)
			} else {
				nt = &Label{
					Name: t,
				}
				relLabels.Set(t, nt)
			}
			nt.Relations = append(nt.Relations, nrel)
		}
	}
	cfg.nEdges = SplitRelations(cfg.Relations)

	for _, k := range relLabels.Keys() {
		nt, _ := relLabels.Get(k)
		cfg.labels = append(cfg.labels, nt.(*Label))
	}

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
		desc, err := cfg.readDescFile(MdPath("_index", []string{}))
		if err != nil {
			return err
		}
		cfg.Desc = desc
	}

	// diagrams
	for _, d := range cfg.Diagrams {
		if d.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_diagram", []string{d.Id()}))
		if err != nil {
			return err
		}
		d.Desc = desc
	}

	// clusters
	for _, c := range cfg.Clusters() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_cluster", []string{c.Id()}))
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
		desc, err := cfg.readDescFile(MdPath("_layer", []string{l.Name}))
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
		desc, err := cfg.readDescFile(MdPath("_node", []string{n.Id()}))
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
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.ClusterComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}
	for _, c := range cfg.NodeComponents() {
		if c.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_component", []string{c.Id()}))
		if err != nil {
			return err
		}
		c.Desc = desc
	}

	// labels
	for _, t := range cfg.labels {
		if t.Desc != "" {
			continue
		}
		desc, err := cfg.readDescFile(MdPath("_label", []string{t.Id()}))
		if err != nil {
			return err
		}
		t.Desc = desc
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
