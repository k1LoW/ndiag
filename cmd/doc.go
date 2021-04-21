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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/output/gviz"
	"github.com/k1LoW/ndiag/output/md"
	"github.com/spf13/cobra"
)

// docCmd represents the doc command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate architecture document",
	Long:  `Generate architecture document.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newConfig()
		if err != nil {
			return err
		}
		format := cfg.Format()

		err = os.MkdirAll(cfg.DocPath, 0755) // #nosec
		if err != nil {
			return err
		}
		if rmDist && cfg.DocPath != "" {
			docs, err := ioutil.ReadDir(cfg.DocPath)
			if err != nil {
				return err
			}
			for _, f := range docs {
				if err := os.RemoveAll(filepath.Join(cfg.DocPath, f.Name())); err != nil {
					return err
				}
			}
		}
		if !force {
			if err := diagExists(cfg); err != nil {
				return err
			}
		}
		err = os.MkdirAll(cfg.DescPath, 0755) // #nosec
		if err != nil {
			return err
		}

		// cleanup empty ndiag.descriptions/*.md
		descs, err := ioutil.ReadDir(cfg.DescPath)
		if err != nil {
			return err
		}
		for _, f := range descs {
			if !f.IsDir() && f.Size() == 0 {
				if err := os.Remove(filepath.Join(cfg.DescPath, f.Name())); err != nil {
					return err
				}
			}
		}

		// views
		for _, v := range cfg.Views {
			cfg, err := newConfig()
			if err != nil {
				return err
			}
			if !cfg.HideViews {
				// generate md
				o := md.New(cfg)
				oldPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("diagram", v.Id()))
				path := filepath.Join(cfg.DocPath, config.MakeMdFilename("view", v.Id()))
				if _, err := os.Stat(oldPath); err == nil {
					if _, err := os.Stat(path); err == nil {
						return fmt.Errorf("old file exists: %s", oldPath)
					}
					if err := os.Rename(oldPath, path); err != nil {
						return err
					}
				}
				file, err := os.Create(path)
				if err != nil {
					return err
				}
				if err := o.OutputView(file, v); err != nil {
					_ = file.Close()
					return err
				}
				_ = file.Close()
			}
			// draw view
			diag := gviz.New(cfg)
			oldPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("diagram", v.Id(), format))
			path := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("view", v.Id(), format))
			if _, err := os.Stat(oldPath); err == nil {
				if _, err := os.Stat(path); err == nil {
					return fmt.Errorf("old diagram file exists: %s", oldPath)
				}
				if err := os.Rename(oldPath, path); err != nil {
					return err
				}
			}
			dFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				return err
			}
			if err := diag.OutputView(dFile, v); err != nil {
				_ = dFile.Close()
				return err
			}
			_ = dFile.Close()
		}

		// layers
		if !cfg.HideLayers {
			for _, l := range cfg.Layers() {
				cfg, err := newConfig()
				if err != nil {
					return err
				}

				// generate md
				o := md.New(cfg)
				mPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("layer", l.Id()))
				file, err := os.Create(mPath)
				if err != nil {
					return err
				}
				if err := o.OutputLayer(file, l); err != nil {
					_ = file.Close()
					return err
				}
				_ = file.Close()

				// draw view
				diag := gviz.New(cfg)
				dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("layer", l.Id(), format))
				dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
				if err != nil {
					return err
				}
				if err := diag.OutputLayer(dFile, l); err != nil {
					_ = dFile.Close()
					return err
				}
				_ = dFile.Close()
			}
		}

		o := md.New(cfg)

		// nodes
		for _, n := range cfg.Nodes {
			cfg, err := newConfig()
			if err != nil {
				return err
			}
			o := md.New(cfg)

			// generate md
			mPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("node", n.Id()))
			file, err := os.Create(mPath)
			if err != nil {
				return err
			}
			if err := o.OutputNode(file, n); err != nil {
				_ = file.Close()
				return err
			}
			_ = file.Close()

			// draw view
			diag := gviz.New(cfg)
			dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("node", n.Id(), format))
			dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				return err
			}
			if err := diag.OutputNode(dFile, n); err != nil {
				_ = dFile.Close()
				return err
			}
			_ = dFile.Close()
		}

		// labels
		if !cfg.HideLabels {
			for k := range cfg.Labels() {
				cfg, err := newConfig()
				if err != nil {
					return err
				}
				rel := cfg.Labels()[k]

				// generate md
				o := md.New(cfg)
				mPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("label", rel.Id()))
				file, err := os.Create(mPath)
				if err != nil {
					return err
				}
				if err := o.OutputLabel(file, rel); err != nil {
					_ = file.Close()
					return err
				}
				_ = file.Close()

				// draw view
				diag := gviz.New(cfg)
				dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("label", rel.Id(), format))
				dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
				if err != nil {
					return err
				}
				if err := diag.OutputLabel(dFile, rel); err != nil {
					_ = dFile.Close()
					return err
				}
				_ = dFile.Close()
			}
		}

		// relations
		for _, r := range cfg.Relations {
			cfg, err := newConfig()
			if err != nil {
				return err
			}

			// draw relation
			diag := gviz.New(cfg)
			dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("relation", r.Id(), format))
			dFile, err := os.OpenFile(dPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
			if err != nil {
				return err
			}
			if err := diag.OutputRelation(dFile, r); err != nil {
				_ = dFile.Close()
				return err
			}
			_ = dFile.Close()
		}

		// top page
		mPath := filepath.Join(cfg.DocPath, "README.md")
		file, err := os.OpenFile(mPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // #nosec
		if err != nil {
			return err
		}
		if err := o.OutputIndex(file); err != nil {
			_ = file.Close()
			return err
		}
		_ = file.Close()

		return nil
	},
}

func diagExists(cfg *config.Config) error {
	format := cfg.Format()
	// views
	for _, d := range cfg.Views {
		mPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("view", d.Id()))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("view", d.Id(), format))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	// layers
	for _, l := range cfg.Layers() {
		mPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("layer", l.Id()))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("layer", l.Id(), format))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	// nodes
	for _, n := range cfg.Nodes {
		mPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("node", n.Id(), format))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("node", n.Id()))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	// labels
	for _, rel := range cfg.Labels() {
		mPath := filepath.Join(cfg.DocPath, config.MakeDiagramFilename("label", rel.Id(), format))
		if _, err := os.Lstat(mPath); err == nil {
			return fmt.Errorf("%s already exist", mPath)
		}
		dPath := filepath.Join(cfg.DocPath, config.MakeMdFilename("label", rel.Id()))
		if _, err := os.Lstat(dPath); err == nil {
			return fmt.Errorf("%s already exist", dPath)
		}
	}

	return nil
}

func init() {
	docCmd.Flags().BoolVarP(&force, "force", "", false, "generate a document without checking for the existence of an existing document")
	docCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	docCmd.Flags().StringSliceVarP(&nodeLists, "nodes", "n", []string{}, "real node list file path")
	docCmd.Flags().BoolVarP(&rmDist, "rm-dist", "", false, "remove all files in the document directory before generating the documents")
	docCmd.Flags().BoolVarP(&hideDetails, "hide-details", "", false, "hide details")
	rootCmd.AddCommand(docCmd)
}
