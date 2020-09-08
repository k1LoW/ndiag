package md

import (
	"io"
	"path/filepath"
	"text/template"

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

	tmpl := template.Must(template.New(d.Name).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Diagram":    d,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
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

	nws := []*config.Network{}
	for _, c := range n.Components {
		for _, nw := range c.Networks {
			nws = append(nws, nw)
		}
	}

	tmpl := template.Must(template.New(n.Id()).Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Node":       n,
		"DiagFormat": m.config.DiagFormat(),
		"DescPath":   rel,
		"Components": n.Components,
		"RealNodes":  n.RealNodes,
		"Networks":   nws,
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
	tmpl := template.Must(template.New("index").Funcs(output.FuncMap).Parse(ts))
	tmplData := map[string]interface{}{
		"Config":     m.config,
		"Diagram":    m.config.PrimaryDiagram(),
		"DiagFormat": m.config.DiagFormat(),
		"Diagrams":   m.config.Diagrams,
		"Layers":     m.config.Layers(),
		"Nodes":      m.config.Nodes,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
