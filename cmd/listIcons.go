package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listIconsCmd = &cobra.Command{
	Use:   "icons",
	Short: "List available icons",
	Long:  `List available icons.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newConfigForIcons()
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Key", "Path"})
		table.SetAutoWrapText(false)
		for _, k := range cfg.IconMap().Keys() {
			i, err := cfg.IconMap().Get(k)
			if err != nil {
				return err
			}
			path := i.Path
			if i.IsGlyph() {
				path = "[embedded icon using github.com/k1LoW/glyph]"
			}
			table.Append([]string{k, path})
		}
		table.Render()
		return nil
	},
}

func init() {
	listIconsCmd.Flags().StringSliceVarP(&configPaths, "config", "c", []string{}, "config file or directory path")
}
