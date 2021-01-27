package md

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/elliotchance/orderedmap"
	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
)

type Md struct {
	config *config.Config
	box    *packr.Box
}

func New(cfg *config.Config) *Md {
	return &Md{
		config: cfg,
		box:    packr.New("md", "./templates"),
	}
}

func (m *Md) OutputDiagram(wr io.Writer, d *config.Diagram) error {
	ts, err := m.box.FindString("diagram.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	layers := []*config.Layer{}
	for _, n := range d.Layers {
		for _, l := range m.config.Layers() {
			if n == l.Name {
				layers = append(layers, l)
			}
		}
	}

	nodes := m.config.Nodes
	labels := m.config.Labels()
	if len(d.Labels) > 0 {
		labels = []*config.Label{}
		for _, t := range d.Labels {
			label, ok := m.config.FindLabel(t)
			if ok != nil {
				return fmt.Errorf("label not found: %s", t)
			}
			labels = append(labels, label)
		}

		nodes, err = m.config.PruneNodesByLabels(nodes, d.Labels)
		if err != nil {
			return err
		}
	}

	tmpl := template.Must(template.New(d.Name).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Diagram":         d,
		"Format":          m.config.Format(),
		"DescPath":        relPath,
		"Layers":          layers,
		"Nodes":           nodes,
		"Labels":          labels,
		"HideLayers":      m.config.HideLayers,
		"HideRealNodes":   m.config.HideRealNodes,
		"HideLabelGroups": m.config.HideLabelGroups,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputLayer(wr io.Writer, l *config.Layer) error {
	ts, err := m.box.FindString("layer.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	clusters, _, _, err := m.config.BuildNestedClusters([]string{l.Name})
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(l.Name).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Layer":         l,
		"Format":        m.config.Format(),
		"DescPath":      relPath,
		"Clusters":      clusters,
		"HideRealNodes": m.config.HideRealNodes,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputNode(wr io.Writer, n *config.Node) error {
	ts, err := m.box.FindString("node.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	labels := []*config.Label{}
	relLabels := orderedmap.NewOrderedMap()
	for _, c := range n.Components {
		for _, e := range c.NEdges {
			for _, ts := range e.Relation.Labels {
				for _, t := range m.config.Labels() {
					if ts == t.Name {
						relLabels.Set(ts, t)
					}
				}
			}
		}
	}
	for _, k := range relLabels.Keys() {
		t, _ := relLabels.Get(k)
		labels = append(labels, t.(*config.Label))
	}

	tmpl := template.Must(template.New(n.Id()).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Node":            n,
		"Format":          m.config.Format(),
		"DescPath":        relPath,
		"Components":      n.Components,
		"RealNodes":       n.RealNodes,
		"Labels":          labels,
		"HideRealNodes":   m.config.HideRealNodes,
		"HideLabelGroups": m.config.HideLabelGroups,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputLabel(wr io.Writer, t *config.Label) error {
	ts, err := m.box.FindString("label.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(t.Id()).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Label":    t,
		"Format":   m.config.Format(),
		"DescPath": relPath,
	}

	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}

	return nil
}

func (m *Md) OutputIndex(wr io.Writer) error {
	ts, err := m.box.FindString("index.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("index").Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Config":   m.config,
		"Diagram":  m.config.PrimaryDiagram(),
		"Format":   m.config.Format(),
		"DescPath": relPath,
		"Diagrams": m.config.Diagrams,
		"Layers":   m.config.Layers(),
		"Nodes":    m.config.Nodes,
		"Labels":   m.config.Labels(),
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
