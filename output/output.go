package output

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/k1LoW/ndiag/config"
)

var unescRep = strings.NewReplacer(fmt.Sprintf("%s%s", config.Esc, config.Sep), config.Sep)

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
	"summary": func(s string) string {
		splitted := strings.Split(s, "\n")
		if len(splitted) <= 1 {
			return splitted[0]
		}
		return fmt.Sprintf("%s ...", splitted[0])
	},
	"imgpath": func(prefix string, vals interface{}, format string) string {
		var strs []string
		switch v := vals.(type) {
		case string:
			strs = []string{v}
		case []string:
			strs = v
		}
		return config.ImagePath(prefix, strs, format)
	},
	"mdpath": func(prefix string, vals interface{}) string {
		var strs []string
		switch v := vals.(type) {
		case string:
			strs = []string{v}
		case []string:
			strs = v
		}
		return config.MdPath(prefix, strs)
	},
}

type Output interface {
	Output(wr io.Writer) error
}
