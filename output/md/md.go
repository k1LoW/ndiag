package md

import (
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

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
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

	tmpl := template.Must(template.New(d.Name).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Diagram":    d,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
		"Layers":     layers,
		"Nodes":      m.config.Nodes,
		"Tags":       m.config.Tags(),
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

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	clusters, _, _, err := m.config.BuildNestedClusters([]string{l.Name})
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(l.Name).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Layer":      l,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
		"Clusters":   clusters,
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

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tags := []*config.Tag{}
	nwTags := orderedmap.NewOrderedMap()
	for _, c := range n.Components {
		for _, e := range c.NEdges {
			for _, ts := range e.Network.Tags {
				for _, t := range m.config.Tags() {
					if ts == t.Name {
						nwTags.Set(ts, t)
					}
				}
			}
		}
	}
	for _, k := range nwTags.Keys() {
		t, _ := nwTags.Get(k)
		tags = append(tags, t.(*config.Tag))
	}

	tmpl := template.Must(template.New(n.Id()).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Node":       n,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
		"Components": n.Components,
		"RealNodes":  n.RealNodes,
		"Tags":       tags,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputTag(wr io.Writer, t *config.Tag) error {
	ts, err := m.box.FindString("network-tag.md.tmpl")
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(t.Id()).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Tag":        t,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
	}

	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}

	return nil
}

func (m *Md) OutputNetwork(wr io.Writer, nw *config.Network) error {
	ts, err := m.box.FindString("network.md.tmpl")
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New(nw.Id()).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Network":    nw,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
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

	rel, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("index").Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Config":     m.config,
		"Diagram":    m.config.PrimaryDiagram(),
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
		"Diagrams":   m.config.Diagrams,
		"Layers":     m.config.Layers(),
		"Nodes":      m.config.Nodes,
		"Tags":       m.config.Tags(),
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
