package gviz

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/antchfx/xmlquery"
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

func (g *Gviz) OutputTag(wr io.Writer, t *config.Tag) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputTag(buf, t); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) OutputRelation(wr io.Writer, rel *config.Relation) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputRelation(buf, rel); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) render(wr io.Writer, b []byte) error {
	format := g.config.Format()
	switch format {
	case "png":
		return g.renderPNG(wr, b)
	case "svg":
		return g.renderSVG(wr, b)
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
}

func (g *Gviz) renderPNG(wr io.Writer, b []byte) (e error) {
	format := g.config.Format()
	_, err := exec.LookPath("dot")
	if err != nil {
		return fmt.Errorf("%v: if the format is png, you need dot command", err)
	}

	tmpIconDir := g.config.TempIconDir()
	if err := os.Mkdir(tmpIconDir, 0750); err != nil {
		return err
	}
	defer os.RemoveAll(tmpIconDir)
	for _, k := range g.config.IconMap().Keys() {
		i, err := g.config.IconMap().Get(k)
		if err != nil {
			return err
		}
		p := filepath.Join(tmpIconDir, fmt.Sprintf("%s.png", k))
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // #nosec
		if err != nil {
			return err
		}
		if err := i.WriteImage(f); err != nil {
			e = f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	// use dot commad
	dotFormatOption := fmt.Sprintf("-T%s", format)
	cmd := exec.Command("dot", dotFormatOption) // #nosec
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if _, err := stdin.Write(b); err != nil {
		_ = stdin.Close()
		return err
	}
	if err := stdin.Close(); err != nil {
		return err
	}
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	if _, err := wr.Write(out); err != nil {
		return err
	}
	return nil
}

func (g *Gviz) renderSVG(wr io.Writer, b []byte) (e error) {
	tmpIconDir := g.config.TempIconDir()
	if err := os.Mkdir(tmpIconDir, 0750); err != nil {
		return err
	}
	defer os.RemoveAll(tmpIconDir)
	for _, k := range g.config.IconMap().Keys() {
		i, err := g.config.IconMap().Get(k)
		if err != nil {
			return err
		}
		p := filepath.Join(tmpIconDir, fmt.Sprintf("%s.svg", k))
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666) // #nosec
		if err != nil {
			return err
		}
		if err := i.Write(f); err != nil {
			e = f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	// use go-graphviz
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

	buf := new(bytes.Buffer)

	if err := gviz.Render(graph, graphviz.Format("svg"), buf); err != nil {
		return err
	}

	doc, err := xmlquery.Parse(buf)
	if err != nil {
		return err
	}

	for _, img := range xmlquery.Find(doc, "//svg/g/g/image") {
		for i, attr := range img.Attr {
			if attr.Name.Space == "xlink" && attr.Name.Local == "href" {
				imgf, err := ioutil.ReadFile(attr.Value)
				if err != nil {
					return err
				}
				imgdoc, err := xmlquery.Parse(bytes.NewReader(imgf))
				if err != nil {
					return err
				}
				s := xmlquery.FindOne(imgdoc, "//svg")
				xmlquery.AddAttr(img, "xlink:href", fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString([]byte(s.OutputXML(true)))))
				img.Attr = append(img.Attr[:i], img.Attr[i+1:]...)
				break
			}
		}
	}

	wr.Write([]byte(doc.OutputXML(false)))

	return nil
}
