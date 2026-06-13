// Package caches detects regenerable caches from common dev tools.
// Items returned here are safe to remove: rebuilding them costs time, not data.
package caches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/item"
)

// Scan runs every cache detector and returns the items found.
// Missing paths are skipped silently; errors are only returned for
// unexpected I/O failures.
func Scan() ([]item.Item, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var items []item.Item

	// Simple path-based caches. All RiskSafe — these regenerate on demand.
	//
	// The KeyBase is used to look up "caches.<key>.name" / ".detail.simple" /
	// ".detail.advanced" in the i18n catalog. Tools without translations fall
	// back to the legacy Name/Detail strings (kept here for safety).
	pathCaches := []struct {
		keyBase, name, tool, rel, detail string
	}{
		{"npm", "npm cache", "npm", ".npm/_cacache", "downloaded packages"},
		{"npm-npx", "npm npx cache", "npm", ".npm/_npx", "one-shot npx executions"},
		{"npm-logs", "npm logs", "npm", ".npm/_logs", "old install logs"},
		{"pnpm", "pnpm store", "pnpm", "Library/pnpm/store", "global content-addressable store"},
		{"yarn", "yarn cache", "yarn", ".yarn/cache", "downloaded packages"},
		{"brew", "Homebrew cache", "brew", "Library/Caches/Homebrew", "downloaded bottles & sources"},
		{"jetbrains", "JetBrains cache", "jetbrains", "Library/Caches/JetBrains", "indexes y logs"},
		{"go", "Go build cache", "go", "Library/Caches/go-build", "compilation cache"},
		{"pip", "pip cache", "pip", "Library/Caches/pip", "wheels & http cache"},
		{"uv", "uv cache", "uv", ".cache/uv", "Python package cache"},
		{"composer", "Composer cache", "composer", "Library/Caches/composer", "PHP packages"},
		{"node-gyp", "node-gyp cache", "node-gyp", "Library/Caches/node-gyp", "native build headers"},
		{"chrome", "Chrome cache", "chrome", "Library/Caches/Google/Chrome", "browser cache"},
		{"firefox", "Firefox cache", "firefox", "Library/Caches/Mozilla", "browser cache"},
		{"xcode-derived", "Xcode DerivedData", "xcode", "Library/Developer/Xcode/DerivedData", "build artifacts"},
		{"xcode-archives", "Xcode Archives", "xcode", "Library/Developer/Xcode/Archives", "old release archives"},
		{"xcode-ios-support", "iOS DeviceSupport", "xcode", "Library/Developer/Xcode/iOS DeviceSupport", "symbol files for old iOS versions"},
		{"xcode-simulator", "CoreSimulator caches", "xcode", "Library/Developer/CoreSimulator/Caches", "simulator caches"},
	}

	for _, pc := range pathCaches {
		path := filepath.Join(home, pc.rel)
		bytes, ok := safeSize(path)
		if !ok {
			continue
		}
		items = append(items, item.Item{
			Name:      pc.name,
			NameKey:   "caches." + pc.keyBase + ".name",
			Tool:      pc.tool,
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryCache,
			Risk:      item.RiskSafe,
			Detail:    pc.detail,
			DetailKey: "caches." + pc.keyBase + ".detail",
		})
	}

	// Cargo caches: registry/{cache,src} and git/checkouts can be removed.
	// We do NOT touch ~/.cargo/bin or registry/index.
	cargoSubs := []struct {
		keyBase, rel, detail string
	}{
		{"cargo-cache", ".cargo/registry/cache", "downloaded crates"},
		{"cargo-src", ".cargo/registry/src", "extracted crate sources"},
		{"cargo-git", ".cargo/git/checkouts", "git dependencies"},
	}
	for _, cs := range cargoSubs {
		path := filepath.Join(home, cs.rel)
		bytes, ok := safeSize(path)
		if !ok {
			continue
		}
		items = append(items, item.Item{
			Name:      "Cargo " + filepath.Base(cs.rel),
			NameKey:   "caches." + cs.keyBase + ".name",
			Tool:      "cargo",
			Path:      path,
			Bytes:     bytes,
			Category:  item.CategoryCache,
			Risk:      item.RiskSafe,
			Detail:    cs.detail,
			DetailKey: "caches." + cs.keyBase + ".detail",
		})
	}

	// JetBrains: detect old IDE versions in Application Support.
	// Each version dir is independent — listing them as separate items lets the user pick.
	if jb, err := jetBrainsOldVersions(home); err == nil {
		items = append(items, jb...)
	}

	// Docker — uses the docker CLI to ask the daemon directly.
	if d, ok := dockerReclaimable(); ok {
		items = append(items, d)
	}

	return items, nil
}

// safeSize returns the size of path in bytes. If path doesn't exist
// or du fails, it returns 0,false so the caller can skip the item.
func safeSize(path string) (int64, bool) {
	if _, err := os.Stat(path); err != nil {
		return 0, false
	}
	bytes, _ := disk.DirSize(path)
	if bytes <= 0 {
		return 0, false
	}
	return bytes, true
}

// jetBrainsOldVersions lists IDE config dirs whose version is older than
// the latest one for that product. Newest version is preserved.
func jetBrainsOldVersions(home string) ([]item.Item, error) {
	root := filepath.Join(home, "Library/Application Support/JetBrains")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil // not installed
	}

	// Group by product prefix (e.g. "PhpStorm" → ["PhpStorm2025.1", "PhpStorm2026.1"]).
	groups := map[string][]string{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		product, version := splitJBVersion(name)
		if version == "" {
			continue // not a versioned product dir (Toolbox, Daemon, etc.)
		}
		groups[product] = append(groups[product], name)
	}

	var items []item.Item
	for product, versions := range groups {
		if len(versions) < 2 {
			continue // only one version installed — keep it
		}
		// Sort lexicographically: "2025.1" < "2025.3" < "2026.1".
		latest := versions[0]
		for _, v := range versions[1:] {
			if v > latest {
				latest = v
			}
		}
		for _, v := range versions {
			if v == latest {
				continue
			}
			path := filepath.Join(root, v)
			bytes, ok := safeSize(path)
			if !ok {
				continue
			}
			items = append(items, item.Item{
				Name:       v,
				Tool:       "jetbrains",
				Path:       path,
				Bytes:      bytes,
				Category:   item.CategoryCache,
				Risk:       item.RiskAskBefore, // settings live here too
				Detail:     fmt.Sprintf("%s old version (latest: %s)", product, latest),
				DetailKey:  "caches.jetbrains-old.detail",
				DetailArgs: []any{product, latest},
			})
		}
	}
	return items, nil
}

// splitJBVersion returns the product name and version suffix.
// "PhpStorm2026.1" → ("PhpStorm", "2026.1"). Returns "", "" if no version.
func splitJBVersion(name string) (string, string) {
	for i, r := range name {
		if r >= '0' && r <= '9' {
			return name[:i], name[i:]
		}
	}
	return "", ""
}

// dockerReclaimable asks the docker daemon for reclaimable space across
// images, build cache, and stopped containers. Volumes are excluded by
// design — they may hold user data and need an explicit opt-in.
func dockerReclaimable() (item.Item, bool) {
	out, err := exec.Command("docker", "system", "df", "--format", "{{.Type}}\t{{.Reclaimable}}").Output()
	if err != nil {
		return item.Item{}, false
	}
	var total int64
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		kind := strings.TrimSpace(parts[0])
		if strings.EqualFold(kind, "Local Volumes") {
			continue
		}
		// Reclaimable column looks like "1.2GB (35%)".
		size := strings.Fields(parts[1])
		if len(size) == 0 {
			continue
		}
		total += parseDockerSize(size[0])
	}
	if total <= 0 {
		return item.Item{}, false
	}
	return item.Item{
		Name:      "Docker reclaimable",
		NameKey:   "caches.docker.name",
		Tool:      "docker",
		Path:      "", // no path: removed via `docker system prune`
		Bytes:     total,
		Category:  item.CategoryCache,
		Risk:      item.RiskSafe,
		Detail:    "images, build cache, stopped containers (volumes excluded)",
		DetailKey: "caches.docker.detail",
	}, true
}

func parseDockerSize(s string) int64 {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0
	}
	multipliers := []struct {
		suffix string
		mult   int64
	}{
		// Order matters: "TB" must be checked before "B".
		{"TB", 1024 * 1024 * 1024 * 1024},
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
		{"B", 1},
	}
	for _, m := range multipliers {
		if strings.HasSuffix(s, m.suffix) {
			num := strings.TrimSuffix(s, m.suffix)
			var val float64
			if _, err := fmt.Sscanf(strings.TrimSpace(num), "%f", &val); err == nil {
				return int64(val * float64(m.mult))
			}
			return 0
		}
	}
	return 0
}
