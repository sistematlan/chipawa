package cmd

import (
	"fmt"
	"os"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/projects"
	"github.com/spf13/cobra"
)

var projectsPath string

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Analiza carpetas de código fuente",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := projects.Scan(projectsPath)
		if err != nil {
			return err
		}

		fmt.Printf("%-35s %-6s %-12s %-30s %s\n", "Proyecto", "Git", "Último commit", "Remote", "Tamaño")
		fmt.Printf("%-35s %-6s %-12s %-30s %s\n", "-------", "---", "-------------", "------", "------")
		for _, p := range list {
			fmt.Printf("%-35s %-6s %-12s %-30s %s\n",
				p.Name, p.HasGitStr(), p.LastCommitLabel(), p.ShortRemote(), disk.FormatBytes(p.Bytes))
		}
		return nil
	},
}

func init() {
	home, _ := os.UserHomeDir()
	projectsCmd.Flags().StringVar(&projectsPath, "path", home+"/sourcecode", "Directorio raíz de proyectos")
	rootCmd.AddCommand(projectsCmd)
}
