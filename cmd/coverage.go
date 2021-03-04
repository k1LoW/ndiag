/*
Copyright © 2021 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/k1LoW/ndiag/coverage"
	"github.com/labstack/gommon/color"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

var cformat string

// coverageCmd represents the coverage command
var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "measure document coverage",
	Long:  `measure document coverage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newConfig()
		if err != nil {
			return err
		}
		cover := coverage.Measure(cfg)
		switch cformat {
		case "json":
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			err := encoder.Encode(cover)
			if err != nil {
				return err
			}
		default:
			max := runewidth.StringWidth("All Elements")
			fmtName := fmt.Sprintf("%%-%ds", max)
			fmt.Printf("%s  %s\n", color.White(fmt.Sprintf(fmtName, "Elements"), color.B), color.White("Coverage", color.B))
			fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "All Elements"), cover.Coverage, cover.Covered, cover.Total)
			if cover.Views != nil {
				fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Views"), cover.Views.Coverage, cover.Views.Covered, cover.Views.Total)
			}
			fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Nodes"), cover.Nodes.Coverage, cover.Nodes.Covered, cover.Nodes.Total)
			fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Components"), cover.Components.Coverage, cover.Components.Covered, cover.Components.Total)
			if cover.Relations != nil {
				fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Relations"), cover.Relations.Coverage, cover.Relations.Covered, cover.Relations.Total)
			}
			fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Layers"), cover.Layers.Coverage, cover.Layers.Covered, cover.Layers.Total)
			if cover.Labels != nil {
				fmt.Printf("%s  %g%% (%d/%d)\n", fmt.Sprintf(fmtName, "Labels"), cover.Labels.Coverage, cover.Labels.Covered, cover.Labels.Total)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(coverageCmd)
	coverageCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	coverageCmd.Flags().StringVarP(&cformat, "format", "t", "", "output format")
}
