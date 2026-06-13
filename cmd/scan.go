package cmd

import (
	"fmt"
	"sort"

	"github.com/sistematlan/mistah/internal/apps"
	"github.com/sistematlan/mistah/internal/caches"
	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/sistematlan/mistah/internal/orphans"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Escaneo completo: disco, apps, caches y datos hu\u00e9rfanos",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== mistah scan ===")
		fmt.Println()

		// 1. Disk usage.
		info, err := disk.Usage("/")
		if err != nil {
			return err
		}
		fmt.Printf("Disco\n  Total: %s  Usado: %s  Libre: %s (%.0f%%)\n\n",
			info.TotalStr, info.UsedStr, info.FreeStr, info.UsedPct)

		// 2. Apps with last-use heuristic.
		appList, err := apps.List()
		if err != nil {
			return err
		}
		unused := 0
		for _, a := range appList {
			if a.DaysSinceUse < 0 || a.DaysSinceUse > 90 {
				unused++
			}
		}
		fmt.Printf("Apps instaladas: %d  (sin uso >90d: %d)\n\n", len(appList), unused)

		// 3. Dev caches — safe to wipe.
		cacheList, err := caches.Scan()
		if err != nil {
			return err
		}
		sort.Slice(cacheList, func(i, j int) bool { return cacheList[i].Bytes > cacheList[j].Bytes })
		fmt.Printf("Caches dev: %s recuperables\n", disk.FormatBytes(item.TotalBytes(cacheList)))
		for _, c := range topN(cacheList, 8) {
			fmt.Printf("  %-30s %s\n", c.Name, disk.FormatBytes(c.Bytes))
		}
		if len(cacheList) > 8 {
			fmt.Printf("  ... y %d m\u00e1s (mistah caches)\n", len(cacheList)-8)
		}
		fmt.Println()

		// 4. Orphans — needs confirmation.
		orphanList, err := orphans.Scan()
		if err != nil {
			return err
		}
		if len(orphanList) > 0 {
			sort.Slice(orphanList, func(i, j int) bool { return orphanList[i].Bytes > orphanList[j].Bytes })
			fmt.Printf("Datos hu\u00e9rfanos: %s (requieren confirmaci\u00f3n)\n",
				disk.FormatBytes(item.TotalBytes(orphanList)))
			for _, o := range orphanList {
				fmt.Printf("  %-30s %s — %s\n", o.Name, disk.FormatBytes(o.Bytes), o.Detail)
			}
			fmt.Println()
		}

		// 5. Grand total.
		grand := item.TotalBytes(cacheList) + item.TotalBytes(orphanList)
		fmt.Printf("Recuperable estimado: %s\n", disk.FormatBytes(grand))
		fmt.Println("Ejecuta `mistah clean --dry-run` para ver qu\u00e9 se borrar\u00eda.")
		return nil
	},
}

// topN returns the first n items, or all if the slice is shorter.
func topN(items []item.Item, n int) []item.Item {
	if len(items) <= n {
		return items
	}
	return items[:n]
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
