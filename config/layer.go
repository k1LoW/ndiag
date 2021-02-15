package config

import (
	"image/color"
	"strings"
)

type Layer struct {
	Name     string
	Desc     string
	Metadata LayerMetadata
}

type LayerMetadata struct {
	Color     color.Color
	FillColor color.Color
	TextColor color.Color
}

func (l *Layer) Id() string {
	return strings.ToLower(l.FullName())
}

func (l *Layer) FullName() string {
	return l.Name
}

func (l *Layer) DescFilename() string {
	return MakeMdFilename("_layer", l.Id())
}

func (l *Layer) String() string {
	return l.Name
}
