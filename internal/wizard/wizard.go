// Package wizard runs the no-args interactive flow: scan, summarize,
// pick a cleaning level, confirm, execute, report.
//
// Audience contract:
//   - The wizard is for people who want mistah to "just work".
//   - It NEVER asks per-item — that's what `mistah clean` is for.
//   - It picks one of three preset levels and shows total impact upfront.
//   - It always offers a final confirm before deleting anything.
//
// Levels:
//
//   Light    — only RiskSafe caches, NO Docker prune. Reversible feel.
//   Standard — Light + Docker prune + JetBrains old versions + Xcode artifacts.
//   Deep     — Standard + orphans + downloads (smart subset).
//
// The level filtering happens here, not in the cleaner, so the cleaner stays
// a pure executor and other entry points (clean cmd) keep their existing
// granular behaviour.
package wizard

import (
	"github.com/sistematlan/mistah/internal/caches"
	"github.com/sistematlan/mistah/internal/downloads"
	"github.com/sistematlan/mistah/internal/item"
	"github.com/sistematlan/mistah/internal/orphans"
)

// Level enumerates the wizard presets.
type Level int

const (
	LevelLight Level = iota
	LevelStandard
	LevelDeep
)

// String returns a stable identifier used for i18n keys and logs.
func (l Level) String() string {
	switch l {
	case LevelLight:
		return "light"
	case LevelStandard:
		return "standard"
	case LevelDeep:
		return "deep"
	default:
		return "unknown"
	}
}

// Inventory is the snapshot the wizard scans before showing the menu.
// Splitting caches/orphans/downloads keeps level-filtering trivial later.
type Inventory struct {
	Caches    []item.Item
	Orphans   []item.Item
	Downloads []item.Item
}

// TotalBytes sums every item the wizard knows about.
func (inv Inventory) TotalBytes() int64 {
	return item.TotalBytes(inv.Caches) +
		item.TotalBytes(inv.Orphans) +
		item.TotalBytes(inv.Downloads)
}

// Scan collects every detector output. It is a thin orchestrator that
// the wizard calls once per session.
func Scan() (Inventory, error) {
	cs, err := caches.Scan()
	if err != nil {
		return Inventory{}, err
	}
	os, err := orphans.Scan()
	if err != nil {
		return Inventory{}, err
	}
	ds, err := downloads.Scan()
	if err != nil {
		return Inventory{}, err
	}
	return Inventory{
		Caches:    cs,
		Orphans:   os,
		Downloads: downloads.AsItems(ds),
	}, nil
}

// PlanFor returns the items that the given level would attempt to remove.
// The same Inventory snapshot is used so totals match what the user saw on screen.
//
// Level rules:
//
//   Light: Caches WITHOUT the Docker entry (no Path) and WITHOUT JetBrains
//          old versions (those touch user settings, RiskAskBefore).
//          Effectively: RiskSafe + Path != "".
//
//   Standard: Light + Docker + JetBrains old versions.
//             Anything with item.CategoryCache, regardless of Risk.
//
//   Deep: Standard + orphans + downloads (smart, large-other excluded already).
func PlanFor(level Level, inv Inventory) []item.Item {
	switch level {
	case LevelLight:
		return filterLight(inv.Caches)
	case LevelStandard:
		return inv.Caches // every cache, including Docker prune and JB old.
	case LevelDeep:
		out := make([]item.Item, 0, len(inv.Caches)+len(inv.Orphans)+len(inv.Downloads))
		out = append(out, inv.Caches...)
		out = append(out, inv.Orphans...)
		out = append(out, inv.Downloads...)
		return out
	default:
		return nil
	}
}

// filterLight keeps only items that are RiskSafe AND have a path.
//
// Why exclude no-path items? The only no-path cache is Docker reclaimable;
// it triggers `docker system prune` which can be surprising to non-devs.
// Light should feel boringly safe — that's its job.
//
// Why also drop JetBrains old versions? They're RiskAskBefore (settings
// folder), so even though they're "caches" they don't fit the Light promise.
func filterLight(items []item.Item) []item.Item {
	out := make([]item.Item, 0, len(items))
	for _, it := range items {
		if it.Risk != item.RiskSafe {
			continue
		}
		if it.Path == "" {
			continue
		}
		out = append(out, it)
	}
	return out
}

// LevelTotals is a convenience for the menu: total bytes per level so the
// user sees "Light (2.3 GB) / Standard (5.6 GB) / Deep (7.1 GB)" upfront.
type LevelTotals struct {
	Light    int64
	Standard int64
	Deep     int64
}

// ComputeTotals returns the byte totals each level would clean.
func ComputeTotals(inv Inventory) LevelTotals {
	return LevelTotals{
		Light:    item.TotalBytes(filterLight(inv.Caches)),
		Standard: item.TotalBytes(inv.Caches),
		Deep:     item.TotalBytes(inv.Caches) + item.TotalBytes(inv.Orphans) + item.TotalBytes(inv.Downloads),
	}
}
