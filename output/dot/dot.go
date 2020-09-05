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
	layers []string
	box    *packr.Box
}

func New(cfg *config.Config, layers []string) *Dot {
	return &Dot{
		config: cfg,
		layers: layers,
		box:    packr.New("dot", "./templates"),
	}
}

func (d *Dot) Output(wr io.Writer) error {
	t := "cluster-diag.dot.tmpl"

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diagram").Funcs(output.FuncMap).Parse(ts))

	clusters, remain, networks, err := d.config.BuildNestedClusters(d.layers)
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
	t := "cluster-diag.dot.tmpl"

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
