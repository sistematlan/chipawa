// Package item defines the common type used across detectors.
// A scannable Item can be a cache, an orphan, a download, or a project —
// anything chipawa might list, summarize, or remove.
package item

// Category groups items by intent so the UI and the cleaner can apply
// different policies (caches are safe to wipe; orphans need confirmation).
type Category string

const (
	CategoryCache    Category = "cache"
	CategoryOrphan   Category = "orphan"
	CategoryDownload Category = "download"
	CategoryProject  Category = "project"
	CategoryApp      Category = "app"
)

// Risk drives the cleaner's confirmation flow.
//
//	Safe       — caches that regenerate, no user data.
//	AskBefore  — touches user data or app state (WhatsApp media, JB caches).
//	Dangerous  — borrar es irreversible y costoso (Docker volumes, code without git).
type Risk string

const (
	RiskSafe      Risk = "safe"
	RiskAskBefore Risk = "ask"
	RiskDangerous Risk = "danger"
)

// Item is a concrete thing detectors return: a path on disk with a size,
// classification, and optional metadata for the UI.
type Item struct {
	// Name shown in the UI (e.g. "npm", "Docker leftover").
	Name string
	// Tool or app of origin (e.g. "docker", "jetbrains"). Empty if N/A.
	Tool string
	// Path is the absolute filesystem path that would be removed.
	// Multiple paths are modelled as multiple Items (keep this 1:1 with rm).
	Path string
	// Bytes used by Path. Set by detectors after measuring.
	Bytes int64
	// Category and Risk classify the item.
	Category Category
	Risk     Risk
	// Detail is a short human note for the UI ("media descargada", "v2025.1 antigua").
	Detail string
}

// Detector reports zero or more items. Detectors must be cheap to call
// and must never mutate the filesystem.
type Detector interface {
	// Name returns a stable identifier ("npm", "docker-leftover").
	Name() string
	// Detect inspects the system and returns the items found.
	// Errors are returned only for unexpected failures; a missing path
	// is not an error — return nil items.
	Detect() ([]Item, error)
}

// TotalBytes sums the size of a slice of items.
func TotalBytes(items []Item) int64 {
	var total int64
	for _, it := range items {
		total += it.Bytes
	}
	return total
}
