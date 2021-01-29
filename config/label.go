package config

import (
	"sort"
	"strings"
)

type Label struct {
	Name      string
	Desc      string
	Relations []*Relation
}

type Labels []*Label

func (l *Label) FullName() string {
	return l.Name
}

func (l *Label) Id() string {
	return strings.ToLower(l.FullName())
}

func (labels Labels) Sort() {
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})
}
