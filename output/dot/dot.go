package dot

import (
	"io"
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
