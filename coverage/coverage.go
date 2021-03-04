package coverage

import (
	"math"

	"github.com/k1LoW/ndiag/config"
)

type Coverage struct {
	Name       string
	Coverage   float64
	Views      *CoverageByElement `json:"views,omitempty"`
	Nodes      *CoverageByElement
	Components *CoverageByElement
	Relations  *CoverageByElement `json:"relations,omitempty"`
	Layers     *CoverageByElement
	Labels     *CoverageByElement `json:"labels,omitempty"`
	Covered    int
	Total      int
}

type CoverageByElement struct {
	Coverage float64
	Covered  int
	Total    int
}

// Measure coverage
func Measure(cfg *config.Config) *Coverage {
	cover := &Coverage{
		Name:       cfg.Name,
		Nodes:      &CoverageByElement{},
		Components: &CoverageByElement{},
		Layers:     &CoverageByElement{},
	}
	// index
	cover.Total += 1
	if cfg.Desc != "" {
		cover.Covered += 1
	}

	// views
	if !cfg.HideViews {
		cover.Views = &CoverageByElement{}
		for _, v := range cfg.Views {
			cover.Views.Total += 1
			if v.Desc != "" {
				cover.Views.Covered += 1
			}
		}
		cover.Views.Coverage = round(float64(cover.Views.Covered) / float64(cover.Views.Total) * 100)
		cover.Total += cover.Views.Total
		cover.Covered += cover.Views.Covered
	}

	// nodes
	for _, n := range cfg.Nodes {
		cover.Nodes.Total += 1
		if n.Desc != "" {
			cover.Nodes.Covered += 1
		}
	}
	cover.Nodes.Coverage = round(float64(cover.Nodes.Covered) / float64(cover.Nodes.Total) * 100)
	cover.Total += cover.Nodes.Total
	cover.Covered += cover.Nodes.Covered

	// components
	for _, c := range cfg.Components() {
		cover.Components.Total += 1
		if c.Desc != "" {
			cover.Components.Covered += 1
		}
	}
	cover.Components.Coverage = round(float64(cover.Components.Covered) / float64(cover.Components.Total) * 100)
	cover.Total += cover.Components.Total
	cover.Covered += cover.Views.Covered

	// relations
	if !(cfg.HideViews && cfg.HideLabels) {
		cover.Relations = &CoverageByElement{}
		for _, r := range cfg.Relations {
			cover.Relations.Total += 1
			if r.Desc != "" {
				cover.Relations.Covered += 1
			}
		}
		cover.Components.Coverage = round(float64(cover.Components.Covered) / float64(cover.Components.Total) * 100)
		cover.Total += cover.Components.Total
		cover.Covered += cover.Components.Covered
	}

	// layers
	for _, l := range cfg.Layers() {
		cover.Layers.Total += 1
		if l.Desc != "" {
			cover.Layers.Covered += 1
		}
	}
	cover.Layers.Coverage = round(float64(cover.Layers.Covered) / float64(cover.Layers.Total) * 100)
	cover.Total += cover.Layers.Total
	cover.Covered += cover.Layers.Covered

	// labels
	if !cfg.HideLabels {
		cover.Labels = &CoverageByElement{}
		for _, l := range cfg.Labels() {
			cover.Labels.Total += 1
			if l.Desc != "" {
				cover.Labels.Covered += 1
			}
		}
		cover.Labels.Coverage = round(float64(cover.Labels.Covered) / float64(cover.Labels.Total) * 100)
		cover.Total += cover.Labels.Total
		cover.Covered += cover.Labels.Covered
	}

	cover.Coverage = round(float64(cover.Covered) / float64(cover.Total) * 100)
	return cover
}

func round(f float64) float64 {
	places := 1
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}
