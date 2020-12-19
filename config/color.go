package config

import (
	"image/color"

	"github.com/muesli/gamut"
)

const DefaultBaseColor = "#4B75B9"
const DefaultTextColor = "#333333"

var defaultColors = []string{
	"#1F91BE",
	"#B2CF3E",
	"#F0BA32",
	"#8858AA",
}

type ColorSet struct {
	Color, FillColor, TextColor color.Color
}

type ColorSets []ColorSet

func (s ColorSets) Get(i int) ColorSet {
	idx := i % len(s)
	return s[idx]
}

// defaultColorSets return color sets for layer line
func defaultColorSets(bc string, tc string) ColorSets {
	csets := ColorSets{}
	for _, c := range defaultColors {
		csets = append(csets, ColorSet{
			Color:     gamut.Hex(c),
			FillColor: gamut.Hex("#FFFFFF00"),
			TextColor: gamut.Hex(tc),
		})
	}
	for _, c := range defaultColors {
		csets = append(csets, ColorSet{
			Color:     gamut.Lighter(gamut.Hex(c), 0.1),
			FillColor: gamut.Hex("#FFFFFF00"),
			TextColor: gamut.Hex(tc),
		})
	}
	return csets
}
