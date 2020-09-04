package output

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/k1LoW/ndiag/config"
)

var unescRep = strings.NewReplacer(fmt.Sprintf("%s%s", config.Esc, config.Sep), config.Sep)

func ImagePath(prefix string, vals []string, format string) string {
	return fmt.Sprintf("%s-%s.%s", prefix, strings.Join(vals, "-"), format)
}

func MdPath(prefix string, vals []string) string {
	return fmt.Sprintf("%s-%s.md", prefix, strings.Join(vals, "-"))
}

var FuncMap = template.FuncMap{
	"id": func(e config.Edge) string {
		return unescRep.Replace(e.Id())
	},
	"fullname": func(e config.Edge) string {
		return unescRep.Replace(e.FullName())
	},
	"unesc": func(s string) string {
		return unescRep.Replace(s)
	},
	"imgpath": func(prefix string, vals interface{}, format string) string {
		var strs []string
		switch v := vals.(type) {
		case string:
			strs = []string{v}
		case []string:
			strs = v
		}
		return ImagePath(prefix, strs, format)
	},
	"mdpath": func(prefix string, vals interface{}) string {
		var strs []string
		switch v := vals.(type) {
		case string:
			strs = []string{v}
		case []string:
			strs = v
		}
		return MdPath(prefix, strs)
	},
}

type Output interface {
	Output(wr io.Writer) error
}
