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
}

func New(cfg *config.Config) *Gviz {
	return &Gviz{
		config: cfg,
		dot:    dot.New(cfg),
	}
}

func (g *Gviz) OutputDiagram(wr io.Writer, d *config.Diagram) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputDiagram(buf, d); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) OutputLayer(wr io.Writer, l *config.Layer) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputLayer(buf, l); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) OutputNode(wr io.Writer, n *config.Node) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputNode(buf, n); err != nil {
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
