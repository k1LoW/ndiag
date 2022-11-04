package dot

import (
	"embed"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/elliotchance/orderedmap"
	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
)

//go:embed templates/*.tmpl
var tmpls embed.FS

type Dot struct {
	config *config.Config
}

func New(cfg *config.Config) *Dot {
	return &Dot{
		config: cfg,
	}
}

func (d *Dot) OutputView(wr io.Writer, v *config.View) error {
	ts, err := tmpls.ReadFile("templates/view.dot.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("view").Funcs(output.Funcs(d.config)).Parse(string(ts)))

	clusters, globalNodes, edges, err := d.config.BuildNestedClusters(v.Layers)
	if err != nil {
		return err
	}
	globalComponents := d.config.GlobalComponents()
	clusters, globalNodes, globalComponents, edges, err = d.config.PruneClustersByLabels(clusters, globalNodes, globalComponents, edges, v.Labels)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, map[string]interface{}{
		"GraphAttrs":       d.config.Graph.Attrs,
		"Clusters":         clusters,
		"GlobalNodes":      globalNodes,
		"GlobalComponents": globalComponents,
		"Edges":            config.MergeEdges(edges),
		"HideUnlinked":     false,
		"HideRealNodes":    d.config.HideRealNodes,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputLayer(wr io.Writer, l *config.Layer) error {
	ts, err := tmpls.ReadFile("templates/view.dot.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("view").Funcs(output.Funcs(d.config)).Parse(string(ts)))

	clusters, globalNodes, edges, err := d.config.BuildNestedClusters([]string{l.Id()})
	if err != nil {
		return err
	}
	filteredEdges := []*config.Edge{}
L:
	for _, e := range edges {
		for _, n := range globalNodes {
			// remove rel with global nodes
			if strings.HasPrefix(e.Src.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
			if strings.HasPrefix(e.Dst.Id(), fmt.Sprintf("%s:", n.Id())) {
				continue L
			}
		}
		// remove rel with global components
		if (e.Src.Node == nil && e.Src.Cluster == nil) || (e.Dst.Node == nil && e.Dst.Cluster == nil) {
			continue L
		}
		filteredEdges = append(filteredEdges, e)
	}

	if err := tmpl.Execute(wr, map[string]interface{}{
		"GraphAttrs":       d.config.Graph.Attrs,
		"Clusters":         clusters,
		"GlobalNodes":      []*config.Node{},
		"GlobalComponents": []*config.Component{},
		"Edges":            config.MergeEdges(filteredEdges),
		"HideUnlinked":     false,
		"HideRealNodes":    d.config.HideRealNodes,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputNode(wr io.Writer, n *config.Node) error {
	ts, err := tmpls.ReadFile("templates/node.dot.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("view").Funcs(output.Funcs(d.config)).Parse(string(ts)))

	clusters := config.Clusters{}
	cIds := orderedmap.NewOrderedMap() // map[string]*config.Cluster{}
	nodes := []*config.Node{n}
	nIds := orderedmap.NewOrderedMap() // map[string]*config.Node{ n.Id(): n }
	nIds.Set(n.Id(), n)
	globalComponents := []*config.Component{}
	gIds := orderedmap.NewOrderedMap() // map[string]*config.Component{}
	edges := []*config.Edge{}

	for _, e := range d.config.Edges() {
		if (e.Src.Node == nil || e.Src.Node.Id() != n.Id()) && (e.Dst.Node == nil || e.Dst.Node.Id() != n.Id()) {
			continue
		}
		switch {
		case e.Src.Node != nil:
			nIds.Set(e.Src.Node.Id(), e.Src.Node)
		case e.Src.Cluster != nil:
			e.Src.Cluster.Nodes = nil
			cIds.Set(e.Src.Cluster.Id(), e.Src.Cluster)
		default:
			gIds.Set(e.Src.Id(), e.Src)
		}
		switch {
		case e.Dst.Node != nil:
			nIds.Set(e.Dst.Node.Id(), e.Dst.Node)
		case e.Dst.Cluster != nil:
			e.Dst.Cluster.Nodes = nil
			cIds.Set(e.Dst.Cluster.Id(), e.Dst.Cluster)
		default:
			gIds.Set(e.Dst.Id(), e.Dst)
		}
		edges = append(edges, e)
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
		"GraphAttrs":       d.config.Graph.Attrs,
		"MainNodeId":       n.Id(),
		"Clusters":         clusters,
		"GlobalNodes":      nodes,
		"GlobalComponents": globalComponents,
		"Edges":            config.MergeEdges(edges),
		"HideRealNodes":    d.config.HideRealNodes,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputLabel(wr io.Writer, l *config.Label) error {
	ts, err := tmpls.ReadFile("templates/view.dot.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("view").Funcs(output.Funcs(d.config)).Parse(string(ts)))

	clusters, globalNodes, edges, err := d.config.BuildNestedClusters([]string{})
	if err != nil {
		return err
	}

	globalComponents := d.config.GlobalComponents()
	clusters, globalNodes, globalComponents, edges, err = d.config.PruneClustersByLabels(clusters, globalNodes, globalComponents, edges, []string{l.Id()})
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, map[string]interface{}{
		"GraphAttrs":       d.config.Graph.Attrs,
		"Clusters":         clusters,
		"GlobalNodes":      globalNodes,
		"GlobalComponents": globalComponents,
		"Edges":            edges,
		"HideUnlinked":     false,
		"HideRealNodes":    d.config.HideRealNodes,
	}); err != nil {
		return err
	}
	return nil
}

func (d *Dot) OutputRelation(wr io.Writer, r *config.Relation) error {
	ts, err := tmpls.ReadFile("templates/view.dot.tmpl")
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("view").Funcs(output.Funcs(d.config)).Parse(string(ts)))

	clusters, globalNodes, _, err := d.config.BuildNestedClusters([]string{})
	if err != nil {
		return err
	}

	globalComponents := d.config.GlobalComponents()
	clusters, globalNodes, globalComponents, edges, err := d.config.PruneClustersByRelations(clusters, globalNodes, globalComponents, config.Relations{r})
	if err != nil {
		return err
	}

	attrs := append(d.config.Graph.Attrs, &config.Attr{
		Key:   "rankdir",
		Value: "LR",
	})

	if err := tmpl.Execute(wr, map[string]interface{}{
		"GraphAttrs":       attrs,
		"Clusters":         clusters,
		"GlobalNodes":      globalNodes,
		"GlobalComponents": globalComponents,
		"Edges":            edges,
		"HideUnlinked":     false,
		"HideRealNodes":    true,
	}); err != nil {
		return err
	}
	return nil
}
