/*
Copyright © 2020 Ken'ichiro Oyama <k1lowxb@gmail.com>

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
	"os"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
	"github.com/k1LoW/ndiag/output/dot"
	"github.com/k1LoW/ndiag/output/gviz"
	"github.com/spf13/cobra"
)

// drawCmd represents the draw command
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "draw diagram",
	Long:  `draw diagram.`,
	Run: func(cmd *cobra.Command, args []string) {
		var o output.Output

		cfg := config.New()
		if err := cfg.LoadConfigFile(configPath); err != nil {
			printFatalln(cmd, err)
		}
		for _, l := range nodeLists {
			if err := cfg.LoadRealNodesFile(l); err != nil {
				printFatalln(cmd, err)
			}
		}
		if err := cfg.Build(); err != nil {
			printFatalln(cmd, err)
		}

		switch format {
		case "svg", "jpg", "png":
			o = gviz.New(cfg, layers, format)
		case "dot":
			o = dot.New(cfg, layers)
		}

		if err := o.Output(os.Stdout); err != nil {
			printFatalln(cmd, err)
		}
	},
}

func init() {
	drawCmd.Flags().StringVarP(&format, "format", "t", config.DefaultDiagFormat, "format")
	drawCmd.Flags().StringSliceVarP(&layers, "layer", "l", []string{}, "layer")
	drawCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	drawCmd.Flags().StringSliceVarP(&nodeLists, "node-list", "n", []string{}, "real node list file path")
	drawCmd.Flags().StringVarP(&out, "out", "", "", "output file path")
	rootCmd.AddCommand(drawCmd)
}
