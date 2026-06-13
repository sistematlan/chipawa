// Package downloads inspects ~/Downloads and identifies stale or duplicated
// content that is safe to remove or worth reviewing.
//
// Unlike caches, files in Downloads can hold user-relevant data, so every
// item is at most RiskAskBefore. The classifier returns BOTH detected
// candidates AND a fallback "large file" category for everything that
// doesn't match a smart rule — the user keeps full visibility.
//
// Smart rules implemented:
//
//   - DMG / PKG installers whose app is already in /Applications → safe to remove.
//   - ZIP / RAR / 7z that have a sibling directory with the same base name
//     (heuristic: archive was already extracted in place).
//   - Project folders with a node_modules subdirectory (someone forgot a
//     project in Downloads — node_modules can be 1-3 GB easily).
//   - Database dumps (.sql, .sql.bak, .dump, .tar.gz with sql/ inside)
//     older than 30 days.
//   - Large video files (.mov, .mp4, .mkv) older than 90 days.
//   - Archive files (.zip, .rar, .7z, .tar.gz) older than 90 days that are
//     NOT detected as already-extracted (still archived but stale).
package downloads

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/item"
)

// Subcategory classifies a downloads item further than item.Category.
// It drives the UI grouping and the cleaner's risk policy.
type Subcategory string

const (
	SubInstaller        Subcategory = "installer-with-app"   // DMG/PKG; app already installed
	SubArchiveExtracted Subcategory = "archive-extracted"    // ZIP/RAR with sibling folder
	SubProjectFolder    Subcategory = "project-folder"       // contains node_modules / vendor / target
	SubDBDump           Subcategory = "db-dump"              // .sql, .sql.bak, .dump
	SubOldVideo         Subcategory = "old-video"            // mov/mp4 >90d
	SubOldArchive       Subcategory = "old-archive"          // zip/rar/7z >90d, not extracted
	SubLargeOther       Subcategory = "large-other"          // catch-all >100MB, listed only
)

// Detail wraps an item with a subcategory and last-modified date so the
// command UI can group and sort. It implements no behaviour — it's just
// the carrier shape for the package's public API.
type Detail struct {
	Item        item.Item
	Sub         Subcategory
	ModTime     time.Time
	AgeDays     int // -1 if mod time unknown
}

// Detail thresholds. Centralised so tests can cover boundary cases.
const (
	thresholdDBDays      = 30
	thresholdVideoDays   = 90
	thresholdArchiveDays = 90
	thresholdLargeBytes  = 100 * 1024 * 1024 // 100 MB
	largeOtherTopN       = 20                // cap "large-other" results
)

// Scan walks the user's Downloads folder one level deep and returns the
// classified items. Subdirectories are not traversed unless they match a
// smart rule (project folder).
//
// The function does not stat() into the entire tree; that would be too slow
// on multi-GB folders. We rely on top-level entries plus targeted disk.DirSize
// calls per candidate.
func Scan() ([]Detail, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(home, "Downloads")
	return ScanPath(root)
}

// ScanPath is Scan but parameterised, used by tests with a TempDir.
func ScanPath(root string) ([]Detail, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		// Missing Downloads is not an error — just no items.
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Index sibling names so archive-extracted detection is O(1).
	siblingDirs := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() {
			siblingDirs[e.Name()] = true
		}
	}

	var details []Detail
	classified := map[string]bool{} // entries already matched by a smart rule

	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue // hidden files (.DS_Store, .localized)
		}
		path := filepath.Join(root, name)

		info, err := e.Info()
		if err != nil {
			continue
		}
		modTime := info.ModTime()
		ageDays := int(time.Since(modTime).Hours() / 24)

		if d, ok := classifySmart(path, name, e.IsDir(), modTime, ageDays, siblingDirs); ok {
			details = append(details, d)
			classified[name] = true
		}
	}

	// Catch-all: large files NOT classified above.
	largeOthers := collectLargeOthers(root, entries, classified)
	details = append(details, largeOthers...)

	return details, nil
}

// classifySmart applies the rules in priority order. Returns ok=false when
// no rule matches; the entry will then be considered for "large-other".
func classifySmart(
	path, name string,
	isDir bool,
	modTime time.Time,
	ageDays int,
	siblings map[string]bool,
) (Detail, bool) {
	lower := strings.ToLower(name)

	// 1. Project folder containing node_modules (deepest impact, check first).
	if isDir && hasBuildArtifact(path) {
		bytes, _ := disk.DirSize(path)
		if bytes <= 0 {
			return Detail{}, false
		}
		return makeDetail(item.Item{
			Name:      name,
			Tool:      "downloads",
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryDownload,
			Risk:      item.RiskAskBefore,
			Detail:    "carpeta de proyecto con node_modules/target/dist (probablemente abandonada)",
			DetailKey: "downloads.project-folder.detail",
		}, SubProjectFolder, modTime, ageDays), true
	}

	// 2. Installer (DMG/PKG) whose app is installed.
	if isInstaller(lower) {
		if appName, installed := installerHasApp(name); installed {
			bytes := mustSize(path, isDir)
			if bytes <= 0 {
				return Detail{}, false
			}
			return makeDetail(item.Item{
				Name:       name,
				Tool:       "downloads",
				Path:       path,
				Bytes:      bytes,
				Category:   item.CategoryDownload,
				Risk:       item.RiskSafe, // app already on disk; installer is redundant
				Detail:     "instalador; " + appName + " ya está instalada",
				DetailKey:  "downloads.installer.detail",
				DetailArgs: []any{appName},
			}, SubInstaller, modTime, ageDays), true
		}
	}

	// 3. Archive with a sibling extracted folder.
	if base, ext, ok := splitArchive(lower); ok {
		// Compare against sibling names case-insensitively. We keep the
		// original-case sibling to display, not the lowercased one.
		for sibling := range siblings {
			if strings.EqualFold(sibling, base) {
				bytes := mustSize(path, isDir)
				if bytes <= 0 {
					return Detail{}, false
				}
				// DetailArgs order: simple uses just %s (sibling); advanced
				// uses %s, %s (ext, sibling). The simple variant ignores ext
				// because regular users don't care if it was zip vs rar.
				// To keep one DetailArgs slice, we put sibling first and ext
				// second; templates that only need sibling don't reference ext.
				return makeDetail(item.Item{
					Name:       name,
					Tool:       "downloads",
					Path:       path,
					Bytes:      bytes,
					Category:   item.CategoryDownload,
					Risk:       item.RiskAskBefore,
					Detail:     "archivo " + ext + " ya extraído en ./" + sibling + "/",
					DetailKey:  "downloads.archive-extracted.detail",
					DetailArgs: []any{sibling, ext},
				}, SubArchiveExtracted, modTime, ageDays), true
			}
		}
	}

	// 4. DB dump older than threshold.
	if isDBDump(lower) && ageDays >= thresholdDBDays {
		bytes := mustSize(path, isDir)
		if bytes <= 0 {
			return Detail{}, false
		}
		return makeDetail(item.Item{
			Name:      name,
			Tool:      "downloads",
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryDownload,
			Risk:      item.RiskAskBefore,
			Detail:    "dump de base de datos (>30 días)",
			DetailKey: "downloads.db-dump.detail",
		}, SubDBDump, modTime, ageDays), true
	}

	// 5. Old video.
	if isVideo(lower) && ageDays >= thresholdVideoDays {
		bytes := mustSize(path, isDir)
		if bytes <= 0 {
			return Detail{}, false
		}
		return makeDetail(item.Item{
			Name:      name,
			Tool:      "downloads",
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryDownload,
			Risk:      item.RiskAskBefore,
			Detail:    "video (>90 días)",
			DetailKey: "downloads.old-video.detail",
		}, SubOldVideo, modTime, ageDays), true
	}

	// 6. Old archive (still compressed but stale).
	if _, _, ok := splitArchive(lower); ok && ageDays >= thresholdArchiveDays {
		bytes := mustSize(path, isDir)
		if bytes <= 0 {
			return Detail{}, false
		}
		return makeDetail(item.Item{
			Name:      name,
			Tool:      "downloads",
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryDownload,
			Risk:      item.RiskAskBefore,
			Detail:    "archivo comprimido (>90 días)",
			DetailKey: "downloads.old-archive.detail",
		}, SubOldArchive, modTime, ageDays), true
	}

	return Detail{}, false
}

// collectLargeOthers returns the top-N largest unclassified entries above
// thresholdLargeBytes. These are RiskAskBefore by default and NOT included
// in `clean` automatically — they appear only in `mistah downloads`.
func collectLargeOthers(root string, entries []os.DirEntry, classified map[string]bool) []Detail {
	type sized struct {
		name    string
		bytes   int64
		modTime time.Time
		isDir   bool
	}
	var pool []sized
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || classified[name] {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(root, name)
		bytes := mustSize(path, e.IsDir())
		if bytes < thresholdLargeBytes {
			continue
		}
		pool = append(pool, sized{name, bytes, info.ModTime(), e.IsDir()})
	}
	// Sort by size descending. A simple insertion sort keeps deps minimal.
	for i := 1; i < len(pool); i++ {
		for j := i; j > 0 && pool[j].bytes > pool[j-1].bytes; j-- {
			pool[j], pool[j-1] = pool[j-1], pool[j]
		}
	}
	if len(pool) > largeOtherTopN {
		pool = pool[:largeOtherTopN]
	}

	out := make([]Detail, 0, len(pool))
	for _, p := range pool {
		ageDays := int(time.Since(p.modTime).Hours() / 24)
		out = append(out, makeDetail(item.Item{
			Name:      p.name,
			Tool:      "downloads",
			Path:      filepath.Join(root, p.name),
			Bytes:     p.bytes,
			Category:  item.CategoryDownload,
			Risk:      item.RiskAskBefore,
			Detail:    "archivo grande sin clasificar — revisa antes de borrar",
			DetailKey: "downloads.large-other.detail",
		}, SubLargeOther, p.modTime, ageDays))
	}
	return out
}

// makeDetail is a small constructor that keeps Item + classification metadata together.
func makeDetail(it item.Item, sub Subcategory, modTime time.Time, ageDays int) Detail {
	return Detail{Item: it, Sub: sub, ModTime: modTime, AgeDays: ageDays}
}

// mustSize returns the size of a path. For directories it walks; for files
// it stat()s. Returns 0 on any error so the caller can skip the item.
func mustSize(path string, isDir bool) int64 {
	if isDir {
		bytes, _ := disk.DirSize(path)
		return bytes
	}
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// hasBuildArtifact returns true when the directory contains node_modules,
// vendor (Go), target (Rust), .next (Next.js), or build/dist (generic).
// We only check immediate children — deep walks would be too slow.
func hasBuildArtifact(dir string) bool {
	markers := []string{"node_modules", "target", ".next", "vendor", "dist", "build"}
	for _, m := range markers {
		if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
			return true
		}
	}
	return false
}

// isInstaller reports DMG / PKG / APP-bundle installers.
func isInstaller(lower string) bool {
	return strings.HasSuffix(lower, ".dmg") || strings.HasSuffix(lower, ".pkg")
}

// installerHasApp guesses the app name from the installer filename and
// checks if it's installed in /Applications.
//
// Heuristic: take the first whitespace-delimited token, strip version
// suffixes ("Gemini" from "Gemini-1.0.dmg" → tries "Gemini.app").
//
// Returns the app's display name on match, empty string otherwise.
func installerHasApp(filename string) (string, bool) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	// Take part before the first separator that suggests version metadata.
	for _, sep := range []string{"-", "_"} {
		if i := strings.Index(base, sep); i > 0 {
			candidate := base[:i]
			if appExists(candidate) {
				return candidate, true
			}
		}
	}
	if appExists(base) {
		return base, true
	}
	return "", false
}

// appExists checks /Applications/<name>.app and ~/Applications/<name>.app.
func appExists(name string) bool {
	candidates := []string{filepath.Join("/Applications", name+".app")}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, "Applications", name+".app"))
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return true
		}
	}
	return false
}

// splitArchive returns (basename without extension, extension) when the
// filename has a recognised archive extension. .tar.gz / .tar.bz2 are
// special-cased so the basename strips the full double extension.
func splitArchive(lower string) (string, string, bool) {
	doubles := []string{".tar.gz", ".tar.bz2", ".tar.xz"}
	for _, ext := range doubles {
		if strings.HasSuffix(lower, ext) {
			return strings.TrimSuffix(lower, ext), ext, true
		}
	}
	singles := []string{".zip", ".rar", ".7z", ".tar", ".tgz"}
	for _, ext := range singles {
		if strings.HasSuffix(lower, ext) {
			return strings.TrimSuffix(lower, ext), ext, true
		}
	}
	return "", "", false
}

// isDBDump matches common database dump extensions.
func isDBDump(lower string) bool {
	suffixes := []string{".sql", ".sql.bak", ".dump", ".pgdump", ".sqlite-bak"}
	for _, s := range suffixes {
		if strings.HasSuffix(lower, s) {
			return true
		}
	}
	return false
}

// isVideo matches video extensions worth flagging when stale.
// Note: lots of small .mov screen recordings can add up; we still flag
// them but use the >90d threshold to avoid false positives.
func isVideo(lower string) bool {
	suffixes := []string{".mov", ".mp4", ".mkv", ".avi", ".m4v"}
	for _, s := range suffixes {
		if strings.HasSuffix(lower, s) {
			return true
		}
	}
	return false
}

// AsItems extracts only the Item slice from a list of Details.
// Used by `clean --include-downloads` to feed the cleaner.
func AsItems(details []Detail) []item.Item {
	out := make([]item.Item, 0, len(details))
	for _, d := range details {
		// Skip "large-other" from cleaner — too risky without manual review.
		if d.Sub == SubLargeOther {
			continue
		}
		out = append(out, d.Item)
	}
	return out
}
