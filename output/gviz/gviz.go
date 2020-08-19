package gviz

import (
	"bytes"
	"io"

	"github.com/goccy/go-graphviz"
	"github.com/k1LoW/ndiag/diag"
	"github.com/k1LoW/ndiag/output/dot"
)

type Gviz struct {
	diag        *diag.Diag
	dot         *dot.Dot
	clusterKeys []string
	format      string
}

func New(d *diag.Diag, clusterKeys []string, format string) *Gviz {
	return &Gviz{
		diag:        d,
		dot:         dot.New(d, clusterKeys),
		clusterKeys: clusterKeys,
		format:      format,
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
	if err := gviz.Render(graph, graphviz.Format(g.format), wr); err != nil {
		return err
	}
	return nil
}
