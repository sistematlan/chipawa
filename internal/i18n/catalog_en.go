package i18n

// english is the source-of-truth catalog. Every key MUST exist here.
// Other languages may be incomplete and fall back to these strings.
//
// Convention for keys:
//
//	<area>.<entity>.<field>[.simple|.advanced]
//
// Examples:
//
//	caches.npm.name             — short label shown in tables
//	caches.npm.detail.simple    — phrasing for non-developers
//	caches.npm.detail.advanced  — phrasing with technical terms
//	risk.safe                   — risk badge label
//	cmd.scan.header             — section header in `mistah scan`
//
// Strings may include fmt directives; callers pass args through T(key, args...).
var english = map[string]string{
	// ---- Risk labels ----
	"risk.safe":      "safe (regenerable cache)",
	"risk.ask":       "ask before (may contain user data)",
	"risk.dangerous": "dangerous (irreversible)",

	// ---- Caches: npm ----
	"caches.npm.name":             "npm cache",
	"caches.npm.detail.advanced":  "downloaded packages",
	"caches.npm.detail.simple":    "Downloaded Node.js packages",

	"caches.npm-npx.name":            "npm npx cache",
	"caches.npm-npx.detail.advanced": "one-shot npx executions",
	"caches.npm-npx.detail.simple":   "Temporary tools run with npx",

	"caches.npm-logs.name":            "npm logs",
	"caches.npm-logs.detail.advanced": "old install logs",
	"caches.npm-logs.detail.simple":   "Old npm install logs",

	// ---- Caches: pnpm ----
	"caches.pnpm.name":            "pnpm store",
	"caches.pnpm.detail.advanced": "global content-addressable store",
	"caches.pnpm.detail.simple":   "Shared cache of pnpm packages",

	// ---- Caches: yarn ----
	"caches.yarn.name":            "yarn cache",
	"caches.yarn.detail.advanced": "downloaded packages",
	"caches.yarn.detail.simple":   "Downloaded Yarn packages",

	// ---- Caches: brew ----
	"caches.brew.name":            "Homebrew cache",
	"caches.brew.detail.advanced": "downloaded bottles & sources",
	"caches.brew.detail.simple":   "Old Homebrew installation files",

	// ---- Caches: jetbrains ----
	"caches.jetbrains.name":            "JetBrains cache",
	"caches.jetbrains.detail.advanced": "indexes y logs",
	"caches.jetbrains.detail.simple":   "Indexes and caches of JetBrains IDEs (PhpStorm, IntelliJ…)",

	"caches.jetbrains-old.detail.simple":   "Old version of %s — current is %s",
	"caches.jetbrains-old.detail.advanced": "%s old version (latest: %s)",

	// ---- Caches: go ----
	"caches.go.name":            "Go build cache",
	"caches.go.detail.advanced": "compilation cache",
	"caches.go.detail.simple":   "Go compilation cache",

	// ---- Caches: pip / uv / composer / node-gyp ----
	"caches.pip.name":            "pip cache",
	"caches.pip.detail.advanced": "wheels & http cache",
	"caches.pip.detail.simple":   "Cache of Python packages",

	"caches.uv.name":            "uv cache",
	"caches.uv.detail.advanced": "Python package cache",
	"caches.uv.detail.simple":   "Cache of the uv Python package manager",

	"caches.composer.name":            "Composer cache",
	"caches.composer.detail.advanced": "PHP packages",
	"caches.composer.detail.simple":   "Cache of PHP packages",

	"caches.node-gyp.name":            "node-gyp cache",
	"caches.node-gyp.detail.advanced": "native build headers",
	"caches.node-gyp.detail.simple":   "Files used to compile native Node.js modules",

	// ---- Caches: browsers ----
	"caches.chrome.name":            "Chrome cache",
	"caches.chrome.detail.advanced": "browser cache",
	"caches.chrome.detail.simple":   "Temporary files of Google Chrome",

	"caches.firefox.name":            "Firefox cache",
	"caches.firefox.detail.advanced": "browser cache",
	"caches.firefox.detail.simple":   "Temporary files of Firefox",

	// ---- Caches: cargo ----
	"caches.cargo-cache.name":            "Cargo registry cache",
	"caches.cargo-cache.detail.advanced": "downloaded crates",
	"caches.cargo-cache.detail.simple":   "Downloaded Rust packages",

	"caches.cargo-src.name":            "Cargo registry sources",
	"caches.cargo-src.detail.advanced": "extracted crate sources",
	"caches.cargo-src.detail.simple":   "Extracted Rust package sources",

	"caches.cargo-git.name":            "Cargo git checkouts",
	"caches.cargo-git.detail.advanced": "git dependencies",
	"caches.cargo-git.detail.simple":   "Rust dependencies from Git",

	// ---- Caches: xcode ----
	"caches.xcode-derived.name":            "Xcode DerivedData",
	"caches.xcode-derived.detail.advanced": "build artifacts",
	"caches.xcode-derived.detail.simple":   "Xcode build files (regenerate when you compile)",

	"caches.xcode-archives.name":            "Xcode Archives",
	"caches.xcode-archives.detail.advanced": "old release archives",
	"caches.xcode-archives.detail.simple":   "Old Xcode app archives",

	"caches.xcode-ios-support.name":            "iOS DeviceSupport",
	"caches.xcode-ios-support.detail.advanced": "symbol files for old iOS versions",
	"caches.xcode-ios-support.detail.simple":   "Symbols for old iOS versions",

	"caches.xcode-simulator.name":            "CoreSimulator caches",
	"caches.xcode-simulator.detail.advanced": "simulator caches",
	"caches.xcode-simulator.detail.simple":   "Cache of iOS simulators",

	// ---- Caches: docker ----
	"caches.docker.name":            "Docker reclaimable",
	"caches.docker.detail.advanced": "images, build cache, stopped containers (volumes excluded)",
	"caches.docker.detail.simple":   "Old Docker images and stopped containers",

	// ---- Orphans ----
	"orphans.docker-leftover.name":            "Docker Desktop leftover",
	"orphans.docker-leftover.detail.advanced": "Docker.app is uninstalled but its container data remains",
	"orphans.docker-leftover.detail.simple":   "Docker data left behind after uninstalling the app",

	"orphans.whatsapp-media.name":            "WhatsApp media",
	"orphans.whatsapp-media.detail.advanced": "downloaded photos/videos/audio (chats not affected)",
	"orphans.whatsapp-media.detail.simple":   "Photos, videos and audio downloaded by WhatsApp (chats are kept)",

	// ---- Downloads subcategories ----
	"downloads.installer.detail.simple":            "Installer; %s is already installed",
	"downloads.installer.detail.advanced":          "installer; %s already installed",
	"downloads.archive-extracted.detail.simple":    "Compressed file already extracted in ./%[1]s/",
	"downloads.archive-extracted.detail.advanced":  "%[2]s archive already extracted in ./%[1]s/",
	"downloads.project-folder.detail.simple":       "Project folder with node_modules (probably abandoned)",
	"downloads.project-folder.detail.advanced":     "carpeta de proyecto con node_modules/target/dist (probablemente abandonada)",
	"downloads.db-dump.detail.simple":              "Old database dump (>30 days)",
	"downloads.db-dump.detail.advanced":            "database dump (>30 days)",
	"downloads.old-video.detail.simple":            "Old video (>90 days)",
	"downloads.old-video.detail.advanced":          "video (>90 days)",
	"downloads.old-archive.detail.simple":          "Old compressed file (>90 days)",
	"downloads.old-archive.detail.advanced":        "compressed archive (>90 days)",
	"downloads.large-other.detail.simple":          "Large file — review before deleting",
	"downloads.large-other.detail.advanced":        "unclassified large file — review before deleting",

	// ---- Categories (subcommand groupings) ----
	"category.cache.simple":   "Temporary files of dev tools",
	"category.cache.advanced": "Developer caches",
	"category.orphan.simple":  "Leftover data from uninstalled apps",
	"category.orphan.advanced": "Orphan data",
	"category.download.simple":  "Files in your Downloads folder",
	"category.download.advanced": "Downloads",

	// ---- Common UI strings ----
	"ui.size":         "size",
	"ui.tool":         "tool",
	"ui.detail":       "detail",
	"ui.note":         "note",
	"ui.path":         "path",
	"ui.risk":         "risk",
	"ui.age":          "age",
	"ui.file":         "file",
	"ui.today":        "today",
	"ui.days-ago":     "%d days ago",
	"ui.never-used":   "never",
	"ui.empty":        "—",
	"ui.total":        "Total",
	"ui.recoverable":  "recoverable",
	"ui.requires-confirmation": "requires confirmation",
	"ui.nothing":      "Nothing to clean. Disk is in order.",
	"ui.dry-run-mode": "Mode: dry-run (nothing will be deleted)",

	// ---- Cleaner prompts ----
	"cleaner.prompt":              "[y/N/v=view/q=quit] ",
	"cleaner.prompt.dangerous":    "TYPE the exact name (%q) or empty to cancel:\n  > ",
	"cleaner.removing":            "removing %s ... ",
	"cleaner.ok":                  "ok",
	"cleaner.failed":              "FAILED: %v",
	"cleaner.dry-prefix":          "[dry-run] would remove: ",
	"cleaner.summary":             "Summary: %d removed, %d skipped, %d failed",
	"cleaner.freed":               "Space freed: %s (of %s planned)",
	"cleaner.plan":                "Plan: %d items, %s recoverable",
	"cleaner.user-quit":           "user quit",
	"cleaner.user-declined":       "user declined",

	// ---- Wizard ----
	"wizard.tagline":             "Clean up your Mac, the open-source way.",
	"wizard.scanning":            "Scanning your disk...",
	"wizard.menu.header":         "What kind of cleanup would you like?",
	"wizard.menu.prompt":         "Choice [1-4]: ",
	"wizard.menu.cancel":         "Cancel — don't delete anything",
	"wizard.cancelled":           "Cancelled. Nothing was deleted.",
	"wizard.thanks":              "Done. Thanks for using mistah.",
	"wizard.level.light.name":    "Light",
	"wizard.level.light.desc":    "Safe dev caches only (npm, brew, pip…)",
	"wizard.level.standard.name": "Standard",
	"wizard.level.standard.desc": "Light + Docker + JetBrains + Xcode artifacts",
	"wizard.level.deep.name":     "Deep",
	"wizard.level.deep.desc":     "Standard + orphan data + Downloads candidates",
	"wizard.confirm.level":       "Level",
	"wizard.confirm.items":       "items",
	"wizard.confirm.prompt":      "Proceed? [y/N] ",
}
