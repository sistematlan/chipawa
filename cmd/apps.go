package cmd

import (
	"fmt"

	"github.com/sistematlan/mistah/internal/apps"
	"github.com/sistematlan/mistah/internal/disk"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Lista apps con fecha de último uso y tamaño",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := apps.List()
		if err != nil {
			return err
		}

		fmt.Printf("%-40s %-10s %s\n", "App", "Tamaño", "Último uso")
		fmt.Printf("%-40s %-10s %s\n", "---", "------", "----------")
		for _, a := range list {
			fmt.Printf("%-40s %-10s %s\n", a.Name, disk.FormatBytes(a.Bytes), a.LastUsedLabel())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
}
