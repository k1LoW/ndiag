package md

import (
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/elliotchance/orderedmap"
	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
)

//go:embed templates/*.tmpl
var tmpls embed.FS

type Md struct {
	config *config.Config
}

func New(cfg *config.Config) *Md {
	return &Md{
		config: cfg,
	}
}

func (m *Md) OutputView(wr io.Writer, v *config.View) error {
	return m.outputView(wr, v, config.TypeView)
}

func (m *Md) OutputLayer(wr io.Writer, l *config.Layer) error {
	ts, err := tmpls.ReadFile("templates/layer.md.tmpl")
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

	tmpl := template.Must(template.New(l.Name).Funcs(output.Funcs(m.config)).Parse(string(ts)))
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
	ts, err := tmpls.ReadFile("templates/node.md.tmpl")
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
		label, ok := l.(*config.Label)
		if !ok {
			continue
		}
		labels = append(labels, label)
	}
	labels.Sort()

	tmpl := template.Must(template.New(n.Id()).Funcs(output.Funcs(m.config)).Parse(string(ts)))
	tmplData := map[string]interface{}{
		"Node":          n,
		"Format":        m.config.Format(),
		"DescPath":      relPath,
		"Components":    n.Components,
		"RealNodes":     n.RealNodes,
		"Labels":        labels,
		"HideRealNodes": m.config.HideRealNodes,
		"HideLabels":    m.config.HideLabels,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) OutputLabel(wr io.Writer, l *config.Label) error {
	v := &config.View{
		Name:   l.Name,
		Desc:   l.Desc,
		Layers: []string{},
		Labels: []string{l.Id()},
	}
	return m.outputView(wr, v, config.TypeLabel)
}

func (m *Md) OutputIndex(wr io.Writer) error {
	ts, err := tmpls.ReadFile("templates/index.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	m.config.Labels().Sort()

	tmpl := template.Must(template.New("index").Funcs(output.Funcs(m.config)).Parse(string(ts)))
	tmplData := map[string]interface{}{
		"Config":   m.config,
		"View":     m.config.PrimaryView(),
		"Format":   m.config.Format(),
		"DescPath": relPath,
		"Views":    m.config.Views,
		"Layers":   m.config.Layers(),
		"Nodes":    m.config.Nodes,
		"Labels":   m.config.Labels(),
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}

func (m *Md) outputView(wr io.Writer, v *config.View, eType config.ElementType) error {
	ts, err := tmpls.ReadFile("templates/view.md.tmpl")
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(filepath.Join("root", m.config.DocPath), filepath.Join("root", m.config.DescPath))
	if err != nil {
		return err
	}

	layers := []*config.Layer{}
	for _, n := range v.Layers {
		for _, l := range m.config.Layers() {
			if n == l.Name {
				layers = append(layers, l)
			}
		}
	}

	nodes := m.config.Nodes
	labels := config.Labels{}
	var relations config.Relations
	if len(v.Labels) > 0 {
		for _, s := range v.Labels {
			label, ok := m.config.FindLabel(s)
			if ok != nil {
				return fmt.Errorf("label not found: %s", s)
			}
			labels = append(labels, label)
		}

		nodes, err = m.config.PruneNodesByLabels(nodes, v.Labels)
		if err != nil {
			return err
		}

		relations = m.config.Relations.FindByLabels(labels)
	} else {
		labels = m.config.Labels()
		relations = m.config.Relations
	}
	labels.Sort()

	hideLabels := m.config.HideLabels
	hideLayers := m.config.HideLayers

	switch eType {
	case config.TypeLabel:
		hideLabels = true
		hideLayers = true
	}

	tmpl := template.Must(template.New(v.Name).Funcs(output.Funcs(m.config)).Parse(string(ts)))
	tmplData := map[string]interface{}{
		"TemplateType":  eType.String(),
		"View":          v,
		"Format":        m.config.Format(),
		"DescPath":      relPath,
		"Nodes":         nodes,
		"Relations":     relations,
		"Layers":        layers,
		"Labels":        labels,
		"HideLayers":    hideLayers,
		"HideRealNodes": m.config.HideRealNodes,
		"HideLabels":    hideLabels,
	}
	if err := tmpl.Execute(wr, tmplData); err != nil {
		return err
	}
	return nil
}
