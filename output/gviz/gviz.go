package gviz

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"

	"github.com/antchfx/xmlquery"
	"github.com/goccy/go-graphviz"
	issvg "github.com/h2non/go-is-svg"
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

func (g *Gviz) OutputView(wr io.Writer, v *config.View) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputView(buf, v); err != nil {
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

func (g *Gviz) OutputLabel(wr io.Writer, l *config.Label) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputLabel(buf, l); err != nil {
		return err
	}
	return g.render(wr, buf.Bytes())
}

func (g *Gviz) OutputRelation(wr io.Writer, r *config.Relation) error {
	buf := &bytes.Buffer{}
	if err := g.dot.OutputRelation(buf, r); err != nil {
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
		return fmt.Errorf("%w: if the format is png, you need dot command", err)
	}

	if err := g.config.IconMap().GeneratePNGGlyphIcons(); err != nil {
		return err
	}
	defer func() {
		if err := g.config.IconMap().RemoveTempIconDir(); err != nil {
			e = err
		}
	}()

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
	if err := g.config.IconMap().GenerateSVGGlyphIcons(); err != nil {
		return err
	}
	defer func() {
		if err := g.config.IconMap().RemoveTempIconDir(); err != nil {
			e = err
		}
	}()

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
				imgf, err := os.ReadFile(attr.Value)
				if err != nil {
					return err
				}
				if issvg.Is(imgf) {
					imgdoc, err := xmlquery.Parse(bytes.NewReader(imgf))
					if err != nil {
						return err
					}
					s := xmlquery.FindOne(imgdoc, "//svg")
					xmlquery.AddAttr(img, "xlink:href", fmt.Sprintf("data:image/svg+xml;base64,%s", base64.StdEncoding.EncodeToString([]byte(s.OutputXML(true)))))
				} else {
					_, format, err := image.DecodeConfig(bytes.NewReader(imgf))
					if err != nil {
						return err
					}
					xmlquery.AddAttr(img, "xlink:href", fmt.Sprintf("data:image/%s;base64,%s", format, base64.StdEncoding.EncodeToString(imgf)))
				}
				img.Attr = append(img.Attr[:i], img.Attr[i+1:]...)
				break
			}
		}
	}

	if _, err := wr.Write([]byte(doc.OutputXML(false))); err != nil {
		return err
	}

	return nil
}
