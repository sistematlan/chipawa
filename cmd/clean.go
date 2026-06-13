package cmd

import (
	"fmt"
	"sort"

	"github.com/sistematlan/mistah/internal/caches"
	"github.com/sistematlan/mistah/internal/cleaner"
	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/downloads"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/sistematlan/mistah/internal/orphans"
	"github.com/spf13/cobra"
)

var (
	cleanDryRun           bool
	cleanYes              bool
	cleanIncludeOrphans   bool
	cleanIncludeDownloads bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Limpieza interactiva con confirmación por ítem",
	Long: `Recorre los caches detectados y, opcionalmente, datos huérfanos y candidatos en Downloads.
Por cada ítem pregunta s/N/v(er)/q(uit). Use --dry-run para previsualizar.

Ejemplos:
  mistah clean --dry-run                # ver qué se borraría
  mistah clean                          # interactivo, solo caches
  mistah clean --include-orphans        # también WhatsApp media, Docker leftover
  mistah clean --include-downloads      # también DMGs ya instalados, ZIPs extraídos, dumps viejos
  mistah clean --yes                    # sin confirmaciones (CI)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Build the candidate list.
		items, err := collectCleanCandidates(cleanIncludeOrphans, cleanIncludeDownloads)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			fmt.Println("Nada que limpiar. Disco en orden.")
			return nil
		}
		// Largest first so the impact is visible up top.
		sort.Slice(items, func(i, j int) bool { return items[i].Bytes > items[j].Bytes })

		// 2. Pick the execution mode.
		mode := cleaner.Interactive
		switch {
		case cleanDryRun:
			mode = cleaner.DryRun
		case cleanYes:
			mode = cleaner.Yes
		}

		// 3. Print a header so the user sees the plan total.
		fmt.Printf("Plan: %d ítems, %s recuperables\n",
			len(items), disk.FormatBytes(item.TotalBytes(items)))
		if mode == cleaner.DryRun {
			fmt.Println("Modo: dry-run (no se borrará nada)")
		}

		// 4. Run.
		var prompter cleaner.Prompter
		if mode == cleaner.Interactive {
			prompter = cleaner.NewTerminalPrompter()
		}
		// Wire the global flag into the cleaner package so prompts use the
		// right phrasing without passing the bool through every layer.
		cleaner.SimpleMode = SimpleMode()
		plan := cleaner.New(items, mode, prompter, cmd.OutOrStdout())
		results := plan.Run()

		// 5. Summary.
		s := cleaner.Summarize(results)
		fmt.Println()
		fmt.Printf("Resumen: %d eliminados, %d omitidos, %d fallidos\n",
			s.Removed, s.Skipped, s.Failed)
		fmt.Printf("Espacio liberado: %s (de %s planificados)\n",
			disk.FormatBytes(s.BytesFreed), disk.FormatBytes(s.BytesPlanned))
		return nil
	},
}

// collectCleanCandidates merges caches and (optionally) orphans/downloads into one list.
// Items with no path AND no specialized remover are filtered out by the
// resolver later, so we don't need to skip them here.
//
// Note: downloads.AsItems() drops the "large-other" subcategory by design —
// those need manual review via `mistah downloads`, not blanket cleaning.
func collectCleanCandidates(includeOrphans, includeDownloads bool) ([]item.Item, error) {
	cs, err := caches.Scan()
	if err != nil {
		return nil, err
	}
	all := cs

	if includeOrphans {
		os, err := orphans.Scan()
		if err != nil {
			return nil, err
		}
		all = append(all, os...)
	}

	if includeDownloads {
		ds, err := downloads.Scan()
		if err != nil {
			return nil, err
		}
		all = append(all, downloads.AsItems(ds)...)
	}

	return all, nil
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Reporta sin borrar")
	cleanCmd.Flags().BoolVar(&cleanYes, "yes", false, "Confirmar todo automáticamente")
	cleanCmd.Flags().BoolVar(&cleanIncludeOrphans, "include-orphans", false,
		"Incluye datos huérfanos (Docker leftover, WhatsApp media, etc.)")
	cleanCmd.Flags().BoolVar(&cleanIncludeDownloads, "include-downloads", false,
		"Incluye candidatos de ~/Downloads (instaladores ya usados, ZIPs extraídos, dumps viejos)")
	rootCmd.AddCommand(cleanCmd)
}
