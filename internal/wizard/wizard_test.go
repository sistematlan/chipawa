package wizard

import (
	"testing"

	"github.com/sistematlan/mistah/internal/item"
)

// fixtureInventory returns a representative inventory:
//   - 2 RiskSafe caches with paths (npm, brew)
//   - 1 RiskSafe cache WITHOUT path (Docker)
//   - 1 RiskAskBefore cache with path (JetBrains old version)
//   - 2 orphans
//   - 2 downloads
func fixtureInventory() Inventory {
	return Inventory{
		Caches: []item.Item{
			{Name: "npm", Path: "/x/npm", Bytes: 100, Risk: item.RiskSafe, Category: item.CategoryCache},
			{Name: "brew", Path: "/x/brew", Bytes: 200, Risk: item.RiskSafe, Category: item.CategoryCache},
			{Name: "docker", Path: "", Bytes: 1000, Risk: item.RiskSafe, Category: item.CategoryCache},
			{Name: "jb-old", Path: "/x/jb", Bytes: 500, Risk: item.RiskAskBefore, Category: item.CategoryCache},
		},
		Orphans: []item.Item{
			{Name: "wpp-media", Path: "/x/wpp", Bytes: 800, Risk: item.RiskAskBefore, Category: item.CategoryOrphan},
			{Name: "docker-leftover", Path: "/x/docker-old", Bytes: 35000, Risk: item.RiskAskBefore, Category: item.CategoryOrphan},
		},
		Downloads: []item.Item{
			{Name: "old.sql", Path: "/x/old.sql", Bytes: 70, Risk: item.RiskAskBefore, Category: item.CategoryDownload},
			{Name: "App.dmg", Path: "/x/App.dmg", Bytes: 80, Risk: item.RiskSafe, Category: item.CategoryDownload},
		},
	}
}

// TestPlanFor_Light: only RiskSafe AND path != "".
// Excludes Docker (no path) and JB old (RiskAskBefore).
func TestPlanFor_Light(t *testing.T) {
	inv := fixtureInventory()
	plan := PlanFor(LevelLight, inv)

	if len(plan) != 2 {
		t.Fatalf("Light should yield 2 items, got %d", len(plan))
	}
	for _, it := range plan {
		if it.Risk != item.RiskSafe {
			t.Errorf("Light included non-safe item: %s (risk=%s)", it.Name, it.Risk)
		}
		if it.Path == "" {
			t.Errorf("Light included path-less item: %s", it.Name)
		}
	}
}

// TestPlanFor_Standard: every cache item, regardless of Risk or Path.
func TestPlanFor_Standard(t *testing.T) {
	inv := fixtureInventory()
	plan := PlanFor(LevelStandard, inv)

	if len(plan) != 4 {
		t.Fatalf("Standard should yield 4 cache items, got %d", len(plan))
	}
	for _, it := range plan {
		if it.Category != item.CategoryCache {
			t.Errorf("Standard leaked non-cache item: %s (cat=%s)", it.Name, it.Category)
		}
	}
}

// TestPlanFor_Deep: caches + orphans + downloads.
func TestPlanFor_Deep(t *testing.T) {
	inv := fixtureInventory()
	plan := PlanFor(LevelDeep, inv)

	want := 4 + 2 + 2
	if len(plan) != want {
		t.Fatalf("Deep should yield %d items, got %d", want, len(plan))
	}

	// All categories should be represented.
	seen := map[item.Category]bool{}
	for _, it := range plan {
		seen[it.Category] = true
	}
	for _, c := range []item.Category{item.CategoryCache, item.CategoryOrphan, item.CategoryDownload} {
		if !seen[c] {
			t.Errorf("Deep missing category %s", c)
		}
	}
}

// TestComputeTotals: Standard >= Light, Deep >= Standard.
func TestComputeTotals(t *testing.T) {
	inv := fixtureInventory()
	tot := ComputeTotals(inv)

	wantLight := int64(100 + 200) // npm + brew
	wantStd := int64(100 + 200 + 1000 + 500)
	wantDeep := wantStd + 800 + 35000 + 70 + 80

	if tot.Light != wantLight {
		t.Errorf("Light total = %d, want %d", tot.Light, wantLight)
	}
	if tot.Standard != wantStd {
		t.Errorf("Standard total = %d, want %d", tot.Standard, wantStd)
	}
	if tot.Deep != wantDeep {
		t.Errorf("Deep total = %d, want %d", tot.Deep, wantDeep)
	}

	// Monotonicity: Light <= Standard <= Deep. The wizard menu must show
	// these in increasing order or users get confused.
	if tot.Light > tot.Standard {
		t.Errorf("Light (%d) should not exceed Standard (%d)", tot.Light, tot.Standard)
	}
	if tot.Standard > tot.Deep {
		t.Errorf("Standard (%d) should not exceed Deep (%d)", tot.Standard, tot.Deep)
	}
}

// TestPlanFor_UnknownLevel returns empty, never panics.
func TestPlanFor_UnknownLevel(t *testing.T) {
	inv := fixtureInventory()
	plan := PlanFor(Level(99), inv)
	if plan != nil {
		t.Errorf("unknown level should yield nil, got %d items", len(plan))
	}
}

// TestLevel_String: stable identifiers used for i18n and logs.
func TestLevel_String(t *testing.T) {
	cases := map[Level]string{
		LevelLight:    "light",
		LevelStandard: "standard",
		LevelDeep:     "deep",
		Level(99):     "unknown",
	}
	for l, want := range cases {
		if got := l.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", l, got, want)
		}
	}
}
