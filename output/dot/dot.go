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
			// remove nw with global nodes
			if strings.HasPrefix(nw.Src.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
			if strings.HasPrefix(nw.Dst.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
		}
		// remove nw with global components
		if (nw.Src.Node == nil && nw.Src.Cluster == nil) || (nw.Dst.Node == nil && nw.Dst.Cluster == nil) {
			continue L
		}
		nws = append(nws, nw)
	}

	if err := tmpl.Execute(wr, map[string]interface{}{
		"Clusters":         clusters,
		"RemainNodes":      []*config.Node{},
		"GlobalComponents": []*config.Component{},
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
		if (nw.Src.Node == nil || nw.Src.Node.Id() != n.Id()) && (nw.Dst.Node == nil || nw.Dst.Node.Id() != n.Id()) {
			continue
		}
		switch {
		case nw.Src.Node != nil:
			nIds[nw.Src.Node.Id()] = nw.Src.Node
		case nw.Src.Cluster != nil:
			nw.Src.Cluster.Nodes = nil
			cIds[nw.Src.Cluster.Id()] = nw.Src.Cluster
		default:
			gIds[nw.Src.Id()] = nw.Src
		}
		switch {
		case nw.Dst.Node != nil:
			nIds[nw.Dst.Node.Id()] = nw.Dst.Node
		case nw.Dst.Cluster != nil:
			nw.Dst.Cluster.Nodes = nil
			cIds[nw.Dst.Cluster.Id()] = nw.Dst.Cluster
		default:
			gIds[nw.Dst.Id()] = nw.Dst
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
