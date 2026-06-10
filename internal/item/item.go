// Package item defines the common type used across detectors.
// A scannable Item can be a cache, an orphan, a download, or a project —
// anything chipawa might list, summarize, or remove.
package item

import "github.com/sistematlan/chipawa/internal/i18n"

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
//
// Localization model:
//   - NameKey + DetailKey point at i18n catalog entries.
//   - DetailArgs are forwarded to fmt.Sprintf when the message has format directives.
//   - Name and Detail are kept as fallback for items that have no catalog entry
//     yet (smooth migration path).
//
// Presenters should call HumanName() / HumanDetail() instead of touching
// Name/Detail directly so they automatically pick the right language and mode.
type Item struct {
	// Name shown in the UI (e.g. "npm", "Docker leftover").
	// Used as fallback when NameKey is empty or its catalog entry missing.
	Name string
	// NameKey is an i18n catalog key like "caches.npm.name".
	NameKey string
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
	// Used as fallback when DetailKey is empty or its catalog entry missing.
	Detail string
	// DetailKey is an i18n catalog key for the detail message. The catalog may
	// have ".simple" and ".advanced" variants — the caller picks via HumanDetail.
	DetailKey string
	// DetailArgs are formatted into the catalog string with fmt.Sprintf.
	// Empty for plain messages.
	DetailArgs []any
}

// HumanName returns the localized short label for the item, falling back
// to Item.Name when no catalog entry exists.
func (it Item) HumanName() string {
	if it.NameKey != "" {
		s := i18n.T(it.NameKey)
		// If T returns the key unchanged it means the catalog miss; fall through.
		if s != it.NameKey {
			return s
		}
	}
	return it.Name
}

// HumanDetail returns the localized detail. simple=true picks the
// non-technical phrasing if available; otherwise the advanced one or
// the raw Detail field.
func (it Item) HumanDetail(simple bool) string {
	if it.DetailKey != "" {
		key := it.DetailKey
		// Try the requested variant first.
		variant := key + ".advanced"
		if simple {
			variant = key + ".simple"
		}
		s := i18n.T(variant, it.DetailArgs...)
		if s != variant { // hit
			return s
		}
		// Fall through to the bare key (no .simple/.advanced suffix).
		s = i18n.T(key, it.DetailArgs...)
		if s != key {
			return s
		}
	}
	return it.Detail
}

// HumanRisk returns the localized risk label.
func (it Item) HumanRisk() string {
	switch it.Risk {
	case RiskSafe:
		return i18n.T("risk.safe")
	case RiskAskBefore:
		return i18n.T("risk.ask")
	case RiskDangerous:
		return i18n.T("risk.dangerous")
	default:
		return string(it.Risk)
	}
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
