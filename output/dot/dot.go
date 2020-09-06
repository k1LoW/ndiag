package dot

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
)

type Dot struct {
	config *config.Config
	box    *packr.Box
}

func New(cfg *config.Config) *Dot {
	return &Dot{
		config: cfg,
		box:    packr.New("dot", "./templates"),
	}
}

func (d *Dot) OutputDiagram(wr io.Writer, diag *config.Diagram) error {
	t := "diagram.dot.tmpl"

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diagram").Funcs(output.FuncMap).Parse(ts))

	clusters, remain, networks, err := d.config.BuildNestedClusters(diag.Layers)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(wr, map[string]interface{}{
		"Clusters":         clusters,
		"RemainNodes":      remain,
		"GlobalComponents": d.config.GlobalComponents(),
		"Networks":         networks,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputLayer(wr io.Writer, l *config.Layer) error {
	t := "diagram.dot.tmpl"

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diagram").Funcs(output.FuncMap).Parse(ts))

	clusters, remain, networks, err := d.config.BuildNestedClusters([]string{l.Name})
	if err != nil {
		return err
	}
	nws := []*config.Network{}
L:
	for _, nw := range networks {
		for _, n := range remain {
			if strings.HasPrefix(nw.Head.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
			if strings.HasPrefix(nw.Tail.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
		}
		nws = append(nws, nw)
	}
	if err := tmpl.Execute(wr, map[string]interface{}{
		"Clusters":         clusters,
		"RemainNodes":      []*config.Node{},
		"GlobalComponents": d.config.GlobalComponents(),
		"Networks":         nws,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputNode(wr io.Writer, n *config.Node) error {
	t := "node.dot.tmpl"

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diagram").Funcs(output.FuncMap).Parse(ts))

	_, _, networks, err := d.config.BuildNestedClusters([]string{})
	if err != nil {
		return err
	}

	clusters := config.Clusters{}
	cIds := map[string]*config.Cluster{}
	nodes := []*config.Node{n}
	nIds := map[string]*config.Node{
		n.Id(): n,
	}
	globalComponents := []*config.Component{}
	gIds := map[string]*config.Component{}

	nws := []*config.Network{}
	for _, nw := range networks {
		if (nw.Head.Node == nil || nw.Head.Node.Id() != n.Id()) && (nw.Tail.Node == nil || nw.Tail.Node.Id() != n.Id()) {
			continue
		}
		switch {
		case nw.Head.Node != nil:
			nIds[nw.Head.Node.Id()] = nw.Head.Node
		case nw.Head.Cluster != nil:
			nw.Head.Cluster.Nodes = nil
			cIds[nw.Head.Cluster.Id()] = nw.Head.Cluster
		default:
			gIds[nw.Head.Id()] = nw.Head
		}
		switch {
		case nw.Tail.Node != nil:
			nIds[nw.Tail.Node.Id()] = nw.Tail.Node
		case nw.Tail.Cluster != nil:
			nw.Tail.Cluster.Nodes = nil
			cIds[nw.Tail.Cluster.Id()] = nw.Tail.Cluster
		default:
			gIds[nw.Tail.Id()] = nw.Tail
		}
		nws = append(nws, nw)
	}
	for _, n := range nIds {
		nodes = append(nodes, n)
	}
	for _, c := range cIds {
		clusters = append(clusters, c)
	}
	for _, c := range gIds {
		globalComponents = append(globalComponents, c)
	}

	if err := tmpl.Execute(wr, map[string]interface{}{
		"MainNodeId":       n.Id(),
		"Clusters":         clusters,
		"RemainNodes":      nodes,
		"GlobalComponents": globalComponents,
		"Networks":         nws,
	}); err != nil {
		return err
	}
	return nil
}
