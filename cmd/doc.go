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
	"fmt"
	"os"
	"path/filepath"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output"
	"github.com/k1LoW/ndiag/output/gviz"
	"github.com/k1LoW/ndiag/output/md"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "doc",
	Long:  `doc.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := newConfig()
		if err != nil {
			printFatalln(cmd, err)
		}
		format := cfg.DiagFormat()

		err = os.MkdirAll(cfg.DocPath, 0755) // #nosec
		if err != nil {
			printFatalln(cmd, err)
		}
		if !force {
			if err := diagExists(cfg); err != nil {
				printFatalln(cmd, err)
			}
		}

		// diagrams
		for _, d := range cfg.Diagrams {
			cfg, err := newConfig()
			if err != nil {
				printFatalln(cmd, err)
			}
			o := md.New(cfg)

			// generate md
			mPath := filepath.Join(cfg.DocPath, output.MdPath("diag", d.Layers))
			file, err := os.Create(mPath)
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := o.OutputDiagram(file, d); err != nil {
				printFatalln(cmd, err)
			}

			// draw diagram
			diag := gviz.New(cfg, d.Layers)
			dPath := filepath.Join(cfg.DocPath, output.ImagePath("diag", d.Layers, format))
			dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := diag.Output(dFile); err != nil {
				printFatalln(cmd, err)
			}
		}

		// layers
		for _, l := range cfg.Layers() {
			cfg, err := newConfig()
			if err != nil {
				printFatalln(cmd, err)
			}
			o := md.New(cfg)

			// generate md
			mPath := filepath.Join(cfg.DocPath, output.MdPath("layer", []string{l}))
			file, err := os.Create(mPath)
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := o.OutputLayer(file, l); err != nil {
				printFatalln(cmd, err)
			}

			// draw diagram
			diag := gviz.New(cfg, []string{l})
			dPath := filepath.Join(cfg.DocPath, output.ImagePath("layer", []string{l}, format))
			if _, err := os.Lstat(dPath); err == nil {
				printFatalln(cmd, fmt.Errorf("%s already exist", dPath))
			}
			dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := diag.Output(dFile); err != nil {
				printFatalln(cmd, err)
			}
		}

		o := md.New(cfg)

		// nodes
		for _, n := range cfg.Nodes {
			mPath := filepath.Join(cfg.DocPath, output.MdPath("node", []string{n.Id()}))
			file, err := os.Create(mPath)
			if err != nil {
				printFatalln(cmd, err)
			}
			if err := o.OutputNode(file, n); err != nil {
				printFatalln(cmd, err)
			}
		}

		// top page
		mPath := filepath.Join(cfg.DocPath, "README.md")
		file, err := os.Create(mPath)
		if err != nil {
			printFatalln(cmd, err)
		}
		if err := o.OutputIndex(file); err != nil {
			printFatalln(cmd, err)
		}

	},
}

func diagExists(cfg *config.Config) error {
	format := cfg.DiagFormat()
	// diagrams
	for _, d := range cfg.Diagrams {
		mPath := filepath.Join(cfg.DocPath, output.MdPath("diag", d.Layers))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, output.ImagePath("diag", d.Layers, format))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	// layers
	for _, l := range cfg.Layers() {
		mPath := filepath.Join(cfg.DocPath, output.MdPath("layer", []string{l}))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, output.ImagePath("layer", []string{l}, format))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	// nodes
	for _, n := range cfg.Nodes {
		mPath := filepath.Join(cfg.DocPath, output.ImagePath("node", []string{n.Id()}, format))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, output.MdPath("node", []string{n.Id()}))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	return nil
}

func newConfig() (*config.Config, error) {
	cfg := config.New()
	if err := cfg.LoadConfigFile(detectConfigPath(configPath)); err != nil {
		return nil, err
	}
	for _, n := range nodeLists {
		if err := cfg.LoadRealNodesFile(n); err != nil {
			return nil, err
		}
	}
	if err := cfg.Build(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func init() {
	docCmd.Flags().BoolVarP(&force, "force", "", false, "force")
	docCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	docCmd.Flags().StringSliceVarP(&nodeLists, "nodes", "n", []string{}, "real node list file path")
	rootCmd.AddCommand(docCmd)
}
