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
}

type Output interface {
	Output(wr io.Writer) error
}
