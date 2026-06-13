package cmd

import (
	"fmt"
	"sort"

	"github.com/sistematlan/mistah/internal/caches"
	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/i18n"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/spf13/cobra"
)

var cachesCmd = &cobra.Command{
	Use:   "caches",
	Short: "Detecta caches de herramientas de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := caches.Scan()
		if err != nil {
			return err
		}
		// Largest first so the user sees impact at the top.
		sort.Slice(list, func(i, j int) bool { return list[i].Bytes > list[j].Bytes })

		simple := SimpleMode()
		// Header row uses i18n keys so labels match the user's language.
		fmt.Printf("%-32s %-10s %-10s %s\n",
			cap1(i18n.T("ui.file")), cap1(i18n.T("ui.tool")), cap1(i18n.T("ui.size")), cap1(i18n.T("ui.detail")))
		fmt.Printf("%-32s %-10s %-10s %s\n", "----", "----", "----", "----")
		for _, c := range list {
			fmt.Printf("%-32s %-10s %-10s %s\n",
				truncate(c.HumanName(), 32), c.Tool, disk.FormatBytes(c.Bytes), c.HumanDetail(simple))
		}
		fmt.Printf("\n%s %s: %s\n",
			i18n.T("ui.total"), i18n.T("ui.recoverable"), disk.FormatBytes(item.TotalBytes(list)))
		return nil
	},
}

// cap1 capitalises the first rune of s. Used for column headers since the
// i18n catalog stores lower-case labels.
func cap1(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 32
	}
	return string(r)
}

// truncate cuts s to n runes appending an ellipsis if it doesn't fit.
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n <= 1 {
		return string(r[:n])
	}
	return string(r[:n-1]) + "…"
}

func init() {
	rootCmd.AddCommand(cachesCmd)
}
