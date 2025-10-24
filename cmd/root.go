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
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/k1LoW/ndiag/config"
	"github.com/k1LoW/ndiag/version"
	"github.com/spf13/cobra"
)

var (
	force       bool
	format      string
	layers      []string
	nodeLists   []string
	configPaths []string
	out         string
	rmDist      bool
	iconPrefix  string
	hideDetails bool
)

var rootCmd = &cobra.Command{
	Use:          "ndiag",
	Short:        `ndiag is a "high-level architecture" viewming/documentation tool`,
	Long:         `ndiag is a "high-level architecture" viewming/documentation tool.`,
	Version:      version.Version,
	SilenceUsage: true,
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	log.SetOutput(io.Discard)
	if env := os.Getenv("DEBUG"); env != "" {
		debug, err := os.Create(fmt.Sprintf("%s.debug", version.Name))
		if err != nil {
			rootCmd.PrintErrln(err)
			os.Exit(1)
		}
		log.SetOutput(debug)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {}

func detectConfigPath(configPath string) string {
	if configPath != "" {
		return configPath
	}
	for _, p := range config.DefaultConfigFilePaths {
		if f, err := os.Stat(p); err == nil && !f.IsDir() {
			return p
		}
	}
	return config.DefaultConfigFilePaths[0]
}

func newConfig() (*config.Config, error) {
	cfg, err := loadConfigFiles(configPaths)
	if err != nil {
		return nil, err
	}
	if len(nodeLists) == 0 {
		cfg.HideRealNodes = true
	} else {
		for _, n := range nodeLists {
			if err := cfg.LoadRealNodesFile(n); err != nil {
				return nil, err
			}
		}
	}
	if err := cfg.Build(); err != nil {
		return nil, err
	}
	if hideDetails {
		if err := cfg.HideDetails(); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func newConfigForIcons() (*config.Config, error) {
	cfg, err := loadConfigFiles(configPaths)
	if err != nil {
		return nil, err
	}
	if err := cfg.BuildForIcons(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadConfigFiles(configPaths []string) (*config.Config, error) {
	cfg := config.New()
	for _, p := range configPaths {
		c := config.New()
		f, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		paths := []string{filepath.Join(p)}
		if f.IsDir() {
			files, err := os.ReadDir(filepath.Join(p))
			if err != nil {
				return nil, err
			}
			paths = []string{}
			for _, file := range files {
				paths = append(paths, filepath.Join(p, file.Name()))
			}
		}
		for _, p := range paths {
			if err := c.LoadConfigFile(detectConfigPath(p)); err != nil {
				return nil, err
			}
			if err := cfg.Merge(c); err != nil {
				return nil, err
			}
		}
	}
	return cfg, nil
}
