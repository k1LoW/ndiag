package gviz

import (
	"bytes"
	"io"

	"github.com/goccy/go-graphviz"
	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output/dot"
)

type Gviz struct {
	config *config.Config
	dot    *dot.Dot
	layers []string
}

func New(cfg *config.Config, layers []string) *Gviz {
	return &Gviz{
		config: cfg,
		dot:    dot.New(cfg, layers),
		layers: layers,
	}
}

func (g *Gviz) Output(wr io.Writer) error {
	buf := &bytes.Buffer{}
	if err := g.dot.Output(buf); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) render(wr io.Writer, b []byte) (e error) {
	format := g.config.DiagFormat()
	gviz := graphviz.New()
	graph, err := graphviz.ParseBytes(b)
	if err != nil {
		return err
	}
	defer func() {
		if err := gviz.Close(); err != nil {
			e = err
		}
		if err := graph.Close(); err != nil {
			e = err
		}
	}()
	if err := gviz.Render(graph, graphviz.Format(format), wr); err != nil {
		return err
	}
	return nil
}
