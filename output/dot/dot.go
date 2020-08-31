package dot

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/k1LoW/ndiag/config"
)

var unescRep = strings.NewReplacer(fmt.Sprintf("%s%s", config.Esc, config.Sep), config.Sep)

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
			return unescRep.Replace(e.Id())
		},
		"fullname": func(e config.Edge) string {
			return unescRep.Replace(e.FullName())
		},
		"unesc": func(s string) string {
			return unescRep.Replace(s)
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
