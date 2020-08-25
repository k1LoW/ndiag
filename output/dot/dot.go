package dot

import (
	"io"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/ndiag/diag"
)

type Dot struct {
	diag        *diag.Diag
	real        bool
	clusterKeys []string
	box         *packr.Box
}

func New(d *diag.Diag, clusterKeys []string) *Dot {
	return &Dot{
		diag:        d,
		clusterKeys: clusterKeys,
		box:         packr.New("dot", "./templates"),
	}
}

func (d *Dot) Output(wr io.Writer) error {
	t := "cluster-diag.dot.tmpl"

	funcMap := template.FuncMap{
		"id": func(e diag.Edge) string {
			return e.Id()
		},
		"fullname": func(e diag.Edge) string {
			return e.FullName()
		},
	}

	ts, err := d.box.FindString(t)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New("diag").Funcs(funcMap).Parse(ts))

	clusters, remain, err := d.diag.BuildNestedClusters(d.clusterKeys)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(wr, map[string]interface{}{
		"Clusters":         clusters,
		"RemainNodes":      remain,
		"GlobalComponents": d.diag.GlobalComponents(),
		"Networks":         d.diag.Networks,
	}); err != nil {
		return err
	}
	return nil
}
