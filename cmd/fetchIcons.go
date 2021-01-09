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
	"github.com/k1LoW/ndiag/icon"
	"github.com/k1LoW/ndiag/icon/aws"
	"github.com/k1LoW/ndiag/icon/gcp"
	"github.com/k1LoW/ndiag/icon/k8s"
	"github.com/spf13/cobra"
)

var fetchIconsCmd = &cobra.Command{
	Use:       "fetch-icons",
	Short:     "Fecth icon set from internet",
	Long:      `Fecth icon set from internet.`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"aws", "gcp", "k8s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		var fetcher icon.Fetcher

		cfg, err := newConfig()
		if err != nil {
			return err
		}

		switch target {
		case "aws":
			fetcher = &aws.AWSIcon{}
		case "gcp":
			fetcher = &gcp.GCPIcon{}
		case "k8s":
			fetcher = &k8s.K8sIcon{}
		}

		if iconPrefix == "" {
			iconPrefix = target
		}
		return fetcher.Fetch(cfg.IconPath, iconPrefix)
	},
}

func init() {
	rootCmd.AddCommand(fetchIconsCmd)
	fetchIconsCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
	fetchIconsCmd.Flags().StringVarP(&iconPrefix, "prefix", "", "", "icon key prefix")
}
