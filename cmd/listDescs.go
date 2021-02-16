package cmd

import (
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listDescsCmd = &cobra.Command{
	Use:   "descs",
	Short: "List description file paths",
	Long:  `List description file paths.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newConfig()
		if err != nil {
			return err
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Type", "ID", "Path"})
		table.SetAutoWrapText(false)
		table.Append([]string{"index (README.md)", "-", filepath.Clean(filepath.Join(cfg.DescPath, "_index.md"))})
		for _, e := range cfg.Elements() {
			table.Append([]string{e.ElementType().String(), e.Id(), filepath.Clean(filepath.Join(cfg.DescPath, e.DescFilename()))})
		}
		table.Render()
		return nil
	},
}

func init() {
	listDescsCmd.Flags().StringVarP(&configPath, "config", "c", "", "config file path")
}
