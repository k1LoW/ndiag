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
	tags := m.config.Tags()
	if len(d.Tags) > 0 {
		tags = []*config.Tag{}
		for _, t := range d.Tags {
			tag, ok := m.config.FindTag(t)
			if ok != nil {
				return fmt.Errorf("tag not found: %s", t)
			}
			tags = append(tags, tag)
		}

		nodes, err = m.config.PruneNodesByTags(nodes, d.Tags)
		if err != nil {
			return err
		}
	}

	tmpl := template.Must(template.New(d.Name).Funcs(output.Funcs(m.config.Dict)).Parse(ts))
	tmplData := map[string]interface{}{
		"Diagram":       d,
		"Format":        m.config.Format(),
		"DescPath":      relPath,
		"Layers":        layers,
		"Nodes":         m.config.Nodes,
		"Tags":          tags,
		"HideLayers":    m.config.HideLayers,
		"HideRealNodes": m.config.HideRealNodes,
		"HideTagGroups": m.config.HideTagGroups,
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

	tmpl := template.Must(template.New(l.Name).Funcs(output.Funcs(m.config.Dict)).Parse(ts))
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

	tags := []*config.Tag{}
	relTags := orderedmap.NewOrderedMap()
	for _, c := range n.Components {
		for _, e := range c.NEdges {
			for _, ts := range e.Relation.Tags {
				for _, t := range m.config.Tags() {
					if ts == t.Name {
						relTags.Set(ts, t)
					}
				}
			}
		}
	}
	for _, k := range relTags.Keys() {
		t, _ := relTags.Get(k)
		tags = append(tags, t.(*config.Tag))
	}

	tmpl := template.Must(template.New(n.Id()).Funcs(output.Funcs(m.config.Dict)).Parse(ts))
	tmplData := map[string]interface{}{
		"Node":          n,
		"Format":        m.config.Format(),
		"DescPath":      relPath,
		"Components":    n.Components,
		"RealNodes":     n.RealNodes,
		"Tags":          tags,
		"HideRealNodes": m.config.HideRealNodes,
		"HideTagGroups": m.config.HideTagGroups,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputTag(wr io.Writer, t *config.Tag) error {
	ts, err := m.box.FindString("tag.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(t.Id()).Funcs(output.Funcs(m.config.Dict)).Parse(ts))
	tmplData := map[string]interface{}{
		"Tag":      t,
		"Format":   m.config.Format(),
		"DescPath": relPath,
	}

	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}

	return nil
}

func (m *Md) OutputRelation(wr io.Writer, rel *config.Relation) error {
	ts, err := m.box.FindString("relation.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(rel.Id()).Funcs(output.Funcs(m.config.Dict)).Parse(ts))
	tmplData := map[string]interface{}{
		"Relation": rel,
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

	tmpl := template.Must(template.New("index").Funcs(output.Funcs(m.config.Dict)).Parse(ts))
	tmplData := map[string]interface{}{
		"Config":   m.config,
		"Diagram":  m.config.PrimaryDiagram(),
		"Format":   m.config.Format(),
		"DescPath": relPath,
		"Diagrams": m.config.Diagrams,
		"Layers":   m.config.Layers(),
		"Nodes":    m.config.Nodes,
		"Tags":     m.config.Tags(),
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
