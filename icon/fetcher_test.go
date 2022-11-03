package icon

import (
	"os"
	"testing"
)

func TestOptimizeSVG(t *testing.T) {
	tests := []struct {
		svg string
	}{
		{"../testdata/crd.svg"},
		{"../testdata/logo_with_border.svg"},
	}
	for _, tt := range tests {
		t.Run(tt.svg, func(t *testing.T) {
			b, err := os.ReadFile(tt.svg)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := OptimizeSVG(b, 100, 100); err != nil {
				t.Fatal(err)
			}
		})
	}
}
