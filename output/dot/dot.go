package dot

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/elliotchance/orderedmap"
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
	cIds := orderedmap.NewOrderedMap() // map[string]*config.Cluster{}
	nodes := []*config.Node{n}
	nIds := orderedmap.NewOrderedMap() // map[string]*config.Node{ n.Id(): n }
	nIds.Set(n.Id(), n)
	globalComponents := []*config.Component{}
	gIds := orderedmap.NewOrderedMap() // map[string]*config.Component{}

	nws := []*config.Network{}
	for _, nw := range networks {
		if (nw.Src.Node == nil || nw.Src.Node.Id() != n.Id()) && (nw.Dst.Node == nil || nw.Dst.Node.Id() != n.Id()) {
			continue
		}
		switch {
		case nw.Src.Node != nil:
			nIds.Set(nw.Src.Node.Id(), nw.Src.Node)
		case nw.Src.Cluster != nil:
			nw.Src.Cluster.Nodes = nil
			cIds.Set(nw.Src.Cluster.Id(), nw.Src.Cluster)
		default:
			gIds.Set(nw.Src.Id(), nw.Src)
		}
		switch {
		case nw.Dst.Node != nil:
			nIds.Set(nw.Dst.Node.Id(), nw.Dst.Node)
		case nw.Dst.Cluster != nil:
			nw.Dst.Cluster.Nodes = nil
			cIds.Set(nw.Dst.Cluster.Id(), nw.Dst.Cluster)
		default:
			gIds.Set(nw.Dst.Id(), nw.Dst)
		}
		nws = append(nws, nw)
	}
	for _, k := range nIds.Keys() {
		n, _ := nIds.Get(k)
		nodes = append(nodes, n.(*config.Node))
	}
	for _, k := range cIds.Keys() {
		c, _ := cIds.Get(k)
		clusters = append(clusters, c.(*config.Cluster))
	}
	for _, k := range gIds.Keys() {
		c, _ := gIds.Get(k)
		globalComponents = append(globalComponents, c.(*config.Component))
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
