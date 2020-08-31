package dot

import (
	"io"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/ndiag/config"
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

	funcMap := template.FuncMap{
		"id": func(e config.Edge) string {
			return e.Id()
		},
		"fullname": func(e config.Edge) string {
			return e.FullName()
		},
	}

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diag").Funcs(funcMap).Parse(ts))

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
