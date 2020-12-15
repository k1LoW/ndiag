package config

import (
	"image/color"

	"github.com/muesli/gamut"
)

const DefaultBaseColor = "#4B75B9"
const DefaultTextColor = "#333333"

type ColorSet struct {
	Color, FillColor, TextColor color.Color
}

// defaultColorSets return color sets for layer line
func defaultColorSets(bc string, tc string) []ColorSet {
	base := gamut.Hex(bc)
	colors := gamut.Tints(base, 8)
	csets := []ColorSet{}
	for _, c := range colors {
		csets = append(csets, ColorSet{
			Color:     c,
			FillColor: gamut.Hex("#FFFFFF00"),
			TextColor: gamut.Hex(tc),
		})
	}
	return csets
}
