package output

import (
	"fmt"
	"image/color"
	"io"
	"strings"
	"text/template"

	"github.com/elliotchance/orderedmap"
	"github.com/k1LoW/ndiag/config"
)

type Output interface {
	OutputView(wr io.Writer, d *config.View) error
}

var unescRep = strings.NewReplacer(fmt.Sprintf("%s%s", config.Esc, config.Sep), config.Sep)
var nl2brRep = strings.NewReplacer("\r\n", "<br>", "\n", "<br>", "\r", "<br>")
var crRep = strings.NewReplacer("\r", "")
var clusterRep = strings.NewReplacer(":", "")

func Funcs(cfg *config.Config) map[string]interface{} {
	return template.FuncMap{
		"trim": func(s string) string {
			return strings.TrimRight(s, "\r\n")
		},
		"nl2br": func(s string) string {
			return nl2brRep.Replace(s)
		},
		"id": func(e config.Element) string {
			return unescRep.Replace(e.Id())
		},
		"fullname": func(e config.Element) string {
			return unescRep.Replace(e.FullName())
		},
		"unesc": func(s string) string {
			return unescRep.Replace(s)
		},
		"node_label": func(n config.Node, hideRealNodes bool) string {
			label := fmt.Sprintf(`%s (%d)`, unescRep.Replace(n.FullName()), len(n.RealNodes))
			if hideRealNodes || len(n.RealNodes) == 0 {
				label = unescRep.Replace(n.Name)
			}
			// Temporarily disable icons to avoid go-graphviz WASM memory issues
			return fmt.Sprintf(`"%s"`, label)
		},
		"cluster_label": func(c config.Cluster) string {
			label := unescRep.Replace(c.FullName())
			// Temporarily disable icons to avoid go-graphviz WASM memory issues
			return fmt.Sprintf(`"%s"`, label)
		},
		"component": func(c config.Component) string {
			bc := cfg.BaseColor
			tc := cfg.TextColor
			label := fmt.Sprintf(`"%s"`, unescRep.Replace(c.Name))
			// Temporarily disable icons to avoid go-graphviz WASM memory issues
			return fmt.Sprintf(`"%s"[label=%s, style="rounded,filled,setlinewidth(3)", color="%s", fillcolor="#FFFFFF", fontcolor="%s" shape=box, fontname="Arial"];`, unescRep.Replace(c.Id()), label, bc, tc)
		},
		"global_component": func(c config.Component) string {
			label := fmt.Sprintf(`"%s"`, unescRep.Replace(c.Name))
			// Temporarily disable icons to avoid go-graphviz WASM memory issues
			return fmt.Sprintf(`"%s"[label=%s, style="rounded,bold", shape=box, fontname="Arial"];`, unescRep.Replace(c.Id()), label)
		},
		"summary": func(s string) string {
			splitted := strings.Split(crRep.Replace(strings.TrimRight(s, "\r\n")), "\n")
			switch {
			case len(splitted) == 0:
				return ""
			case len(splitted) == 1:
				return strings.TrimLeft(splitted[0], "# ")
			case len(splitted) == 2 && splitted[1] == "":
				return strings.TrimLeft(splitted[0], "# ")
			default:
				return fmt.Sprintf("%s ...", strings.TrimLeft(splitted[0], "# "))
			}
		},
		"diagpath": func(prefix string, e config.Element, format string) string {
			return config.MakeDiagramFilename(prefix, e.Id(), format)
		},
		"mdpath": func(prefix string, e config.Element) string {
			return config.MakeMdFilename(prefix, e.Id())
		},
		"componentlink": componentLink,
		"rellink":       relLink,
		"fromlinks": func(edges []*config.Edge, base *config.Component) string {
			links := []string{}
			for _, e := range edges {
				if e.Src.Id() != base.Id() {
					links = append(links, componentLink(e.Src))
				}
			}
			return strings.Join(unique(links), " / ")
		},
		"tolinks": func(edges []*config.Edge, base *config.Component) string {
			links := []string{}
			for _, e := range edges {
				if e.Dst.Id() != base.Id() {
					links = append(links, componentLink(e.Dst))
				}
			}
			return strings.Join(unique(links), " / ")
		},
		"attrs": func(attrs []*config.Attr) string {
			if len(attrs) == 0 {
				return ""
			}
			var out string
			for _, a := range attrs {
				out = fmt.Sprintf("%s, %s=\"%s\"", out, a.Key, a.Value)
			}
			return out
		},
		"dict": func(v ...interface{}) map[string]interface{} {
			dict := map[string]interface{}{}
			length := len(v)
			for i := 0; i < length; i += 2 {
				key, ok := v[i].(string)
				if !ok {
					continue
				}
				dict[key] = v[i+1]
			}
			return dict
		},
		"is_linked": func(c *config.Component, edges []*config.Edge) bool {
			for _, e := range edges {
				if c.Id() == e.Src.Id() || c.Id() == e.Dst.Id() {
					return true
				}
			}
			return false
		},
		"lookup": func(text string) string {
			return cfg.Dict.Lookup(text)
		},
		"colorhex": colorToHex,
	}
}

// componentLink.
func componentLink(c *config.Component) string {
	switch {
	case c.Node != nil:
		return fmt.Sprintf("[%s](%s)", c.Id(), config.MakeMdFilename("node", c.Node.Id()))
	case c.Cluster != nil:
		return fmt.Sprintf("[%s](%s#%s)", c.Id(), config.MakeMdFilename("layer", c.Cluster.Layer.Id()), clusterRep.Replace(c.Cluster.Id()))
	default:
		return c.Id()
	}
}

func relLink(rel *config.Relation) string {
	cIds := []string{}
	for _, r := range rel.Components {
		cIds = append(cIds, r.FullName())
	}
	return fmt.Sprintf("[%s](%s)", strings.Join(cIds, " -> "), config.MakeMdFilename("relation", rel.Id()))
}

func unique(in []string) []string {
	m := orderedmap.NewOrderedMap()
	for _, s := range in {
		m.Set(s, s)
	}
	u := []string{}
	for _, k := range m.Keys() {
		s, _ := m.Get(k)
		str, ok := s.(string)
		if !ok {
			continue
		}
		u = append(u, str)
	}
	return u
}

func colorToHex(c color.Color) string {
	rgba, ok := color.RGBAModel.Convert(c).(color.RGBA)
	if !ok {
		return "#00000000"
	}
	return fmt.Sprintf("#%.2x%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B, rgba.A)
}
