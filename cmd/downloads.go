package cmd

import (
	"fmt"
	"sort"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/downloads"
	"github.com/sistematlan/mistah/internal/i18n"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/spf13/cobra"
)

var downloadsCmd = &cobra.Command{
	Use:   "downloads",
	Short: "Analiza ~/Downloads y detecta archivos candidatos a borrar",
	Long: `Inspecciona ~/Downloads agrupando por tipo:
  • installer-with-app: DMG/PKG cuya app ya está instalada
  • archive-extracted : ZIP/RAR/7z con carpeta extraída al lado
  • project-folder    : carpeta con node_modules/target/.next dentro
  • db-dump           : .sql, .sql.bak, .dump (>30 días)
  • old-video         : .mov/.mp4 (>90 días)
  • old-archive       : .zip/.rar (>90 días, no extraído)
  • large-other       : archivos grandes sin clasificar (>100 MB)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		details, err := downloads.Scan()
		if err != nil {
			return err
		}
		if len(details) == 0 {
			fmt.Println(i18n.T("ui.nothing"))
			return nil
		}
		simple := SimpleMode()

		// Group by subcategory while keeping each group sorted by size desc.
		groups := groupBySubcategory(details)
		order := []downloads.Subcategory{
			downloads.SubProjectFolder,
			downloads.SubInstaller,
			downloads.SubArchiveExtracted,
			downloads.SubDBDump,
			downloads.SubOldVideo,
			downloads.SubOldArchive,
			downloads.SubLargeOther,
		}

		var grandTotal int64
		for _, sub := range order {
			items := groups[sub]
			if len(items) == 0 {
				continue
			}
			sort.Slice(items, func(i, j int) bool { return items[i].Item.Bytes > items[j].Item.Bytes })

			total := sumBytes(items)
			fmt.Printf("\n[%s] %d · %s\n", sub, len(items), disk.FormatBytes(total))
			fmt.Printf("%-50s %-10s %-8s %s\n",
				cap1(i18n.T("ui.file")), cap1(i18n.T("ui.size")),
				cap1(i18n.T("ui.age")), cap1(i18n.T("ui.note")))
			fmt.Printf("%-50s %-10s %-8s %s\n", "----", "----", "----", "----")
			for _, d := range items {
				fmt.Printf("%-50s %-10s %-8s %s\n",
					truncate(d.Item.Name, 50),
					disk.FormatBytes(d.Item.Bytes),
					ageLabel(d.AgeDays),
					d.Item.HumanDetail(simple))
			}
			grandTotal += total
		}

		fmt.Printf("\n%s: %s\n", i18n.T("ui.total"), disk.FormatBytes(grandTotal))
		return nil
	},
}

// groupBySubcategory buckets details so the command prints each section once.
func groupBySubcategory(details []downloads.Detail) map[downloads.Subcategory][]downloads.Detail {
	m := map[downloads.Subcategory][]downloads.Detail{}
	for _, d := range details {
		m[d.Sub] = append(m[d.Sub], d)
	}
	return m
}

func sumBytes(details []downloads.Detail) int64 {
	var total int64
	for _, d := range details {
		total += d.Item.Bytes
	}
	return total
}

// ageLabel renders a compact age string for the table.
//
//	-1 days     → "—"
//	0 days      → "hoy"
//	<60 days    → "Nd"
//	<24 months  → "Nm"
//	otherwise   → "Ny"
func ageLabel(days int) string {
	switch {
	case days < 0:
		return "—"
	case days == 0:
		return i18n.T("ui.today")
	case days < 60:
		return fmt.Sprintf("%dd", days)
	case days < 730:
		return fmt.Sprintf("%dm", days/30)
	default:
		return fmt.Sprintf("%dy", days/365)
	}
}

// avoid unused import warning if item not referenced directly elsewhere
var _ = item.CategoryDownload

func init() {
	rootCmd.AddCommand(downloadsCmd)
}
