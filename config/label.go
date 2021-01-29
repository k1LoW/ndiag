package config

import "strings"

type Label struct {
	Name      string
	Desc      string
	Relations []*Relation
}

func (t *Label) FullName() string {
	return t.Name
}

func (t *Label) Id() string {
	return strings.ToLower(t.FullName())
}
