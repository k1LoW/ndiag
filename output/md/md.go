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

func (m *Md) OutputView(wr io.Writer, d *config.View) error {
	ts, err := m.box.FindString("view.md.tmpl")
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
	labels := config.Labels{}
	if len(d.Labels) > 0 {
		labels = config.Labels{}
		for _, s := range d.Labels {
			label, ok := m.config.FindLabel(s)
			if ok != nil {
				return fmt.Errorf("label not found: %s", s)
			}
			labels = append(labels, label)
		}

		nodes, err = m.config.PruneNodesByLabels(nodes, d.Labels)
		if err != nil {
			return err
		}
	} else {
		labels = m.config.Labels()
	}
	labels.Sort()

	tmpl := template.Must(template.New(d.Name).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"View":         d,
		"Format":          m.config.Format(),
		"DescPath":        relPath,
		"Layers":          layers,
		"Nodes":           nodes,
		"Labels":          labels,
		"HideLayers":      m.config.HideLayers,
		"HideRealNodes":   m.config.HideRealNodes,
		"HideLabels": m.config.HideLabels,
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

	labels := config.Labels{}
	relLabels := orderedmap.NewOrderedMap()
	for _, c := range n.Components {
		for _, l := range c.Labels {
			relLabels.Set(l.Id(), l)
		}
		for _, e := range c.Edges {
			for _, l := range e.Relation.Labels {
				relLabels.Set(l.Id(), l)
			}
		}
	}
	for _, k := range relLabels.Keys() {
		l, _ := relLabels.Get(k)
		labels = append(labels, l.(*config.Label))
	}
	labels.Sort()

	tmpl := template.Must(template.New(n.Id()).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Node":            n,
		"Format":          m.config.Format(),
		"DescPath":        relPath,
		"Components":      n.Components,
		"RealNodes":       n.RealNodes,
		"Labels":          labels,
		"HideRealNodes":   m.config.HideRealNodes,
		"HideLabels": m.config.HideLabels,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputLabel(wr io.Writer, l *config.Label) error {
	ts, err := m.box.FindString("label.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(l.Id()).Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Label":    l,
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

	m.config.Labels().Sort()

	tmpl := template.Must(template.New("index").Funcs(output.Funcs(m.config)).Parse(ts))
	tmplData := map[string]interface{}{
		"Config":   m.config,
		"View":  m.config.PrimaryView(),
		"Format":   m.config.Format(),
		"DescPath": relPath,
		"Views": m.config.Views,
		"Layers":   m.config.Layers(),
		"Nodes":    m.config.Nodes,
		"Labels":   m.config.Labels(),
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
