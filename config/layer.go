package config

import "image/color"

type Layer struct {
	Name      string
	Desc      string
	Meatadata LayerMetadata
}

type LayerMetadata struct {
	Color     color.Color
	FillColor color.Color
	TextColor color.Color
}

func (l Layer) String() string {
	return l.Name
}
