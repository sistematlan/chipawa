package cmd

import (
	"fmt"
	"sort"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/i18n"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/sistematlan/mistah/internal/orphans"
	"github.com/spf13/cobra"
)

var orphansCmd = &cobra.Command{
	Use:   "orphans",
	Short: "Datos huérfanos: apps desinstaladas, media inflada",
	Long: "Detecta directorios grandes que sobreviven cuando una app se desinstala\n" +
		"o que crecen sin límite (p.ej. media de WhatsApp). Estos ítems no son\n" +
		"caches: pueden contener datos del usuario y requieren confirmación.",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := orphans.Scan()
		if err != nil {
			return err
		}
		if len(list) == 0 {
			fmt.Println(i18n.T("ui.nothing"))
			return nil
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Bytes > list[j].Bytes })
		simple := SimpleMode()

		fmt.Printf("%-32s %-10s %-10s %s\n",
			cap1(i18n.T("ui.file")), cap1(i18n.T("ui.tool")),
			cap1(i18n.T("ui.size")), cap1(i18n.T("ui.detail")))
		fmt.Printf("%-32s %-10s %-10s %s\n", "----", "----", "----", "----")
		for _, o := range list {
			fmt.Printf("%-32s %-10s %-10s %s\n",
				truncate(o.HumanName(), 32), o.Tool,
				disk.FormatBytes(o.Bytes), o.HumanDetail(simple))
		}
		fmt.Printf("\n%s: %s — %s\n",
			i18n.T("ui.total"), disk.FormatBytes(item.TotalBytes(list)),
			i18n.T("ui.requires-confirmation"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orphansCmd)
}
