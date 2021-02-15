package config

import (
	"sort"
	"strings"
)

type Label struct {
	Name string
	Desc string
}

type Labels []*Label

func (l *Label) FullName() string {
	return l.Name
}

func (l *Label) Id() string {
	return strings.ToLower(l.FullName())
}

func (l *Label) DescFilename() string {
	return MakeMdFilename("_label", l.Id())
}

func (labels Labels) Sort() {
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})
}

func (labels Labels) Subtract(labels2 Labels) Labels {
	s := Labels{}
	for _, l := range labels {
		for _, l2 := range labels2 {
			if l.Id() == l2.Id() {
				s = append(s, l)
			}
		}
	}
	return s
}
