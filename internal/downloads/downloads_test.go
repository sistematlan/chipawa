package downloads

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sistematlan/chipawa/internal/item"
)

// makeFile writes a file of N bytes at path with given mod time.
func makeFile(t *testing.T, path string, size int, mod time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
}

// makeDir creates a directory with mod time.
func makeDir(t *testing.T, path string, mod time.Time) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
}

func TestSplitArchive(t *testing.T) {
	cases := []struct {
		in       string
		base     string
		ext      string
		ok       bool
	}{
		{"foo.zip", "foo", ".zip", true},
		{"bar.tar.gz", "bar", ".tar.gz", true},
		{"baz.7z", "baz", ".7z", true},
		{"hello.tar", "hello", ".tar", true},
		{"world.txt", "", "", false},
		{"image.png", "", "", false},
	}
	for _, tc := range cases {
		base, ext, ok := splitArchive(tc.in)
		if ok != tc.ok || base != tc.base || ext != tc.ext {
			t.Errorf("splitArchive(%q) = (%q,%q,%v); want (%q,%q,%v)",
				tc.in, base, ext, ok, tc.base, tc.ext, tc.ok)
		}
	}
}

func TestIsDBDump(t *testing.T) {
	yes := []string{"dump.sql", "backup.sql.bak", "data.dump", "old.pgdump"}
	no := []string{"image.png", "archive.zip", "video.mov"}
	for _, n := range yes {
		if !isDBDump(n) {
			t.Errorf("isDBDump(%q) = false; want true", n)
		}
	}
	for _, n := range no {
		if isDBDump(n) {
			t.Errorf("isDBDump(%q) = true; want false", n)
		}
	}
}

func TestHasBuildArtifact(t *testing.T) {
	tmp := t.TempDir()
	if hasBuildArtifact(tmp) {
		t.Errorf("empty dir should NOT have build artifact")
	}
	if err := os.MkdirAll(filepath.Join(tmp, "node_modules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !hasBuildArtifact(tmp) {
		t.Errorf("dir with node_modules SHOULD have build artifact")
	}
}

// TestScanPath_DetectsExtractedArchive: a .zip with a sibling folder of
// the same base name → SubArchiveExtracted.
func TestScanPath_DetectsExtractedArchive(t *testing.T) {
	tmp := t.TempDir()
	old := time.Now().AddDate(0, 0, -10) // 10 days ago, NOT old-archive

	makeFile(t, filepath.Join(tmp, "Project.zip"), 200*1024*1024, old) // 200 MB to surpass large threshold? we want explicit test
	makeDir(t, filepath.Join(tmp, "Project"), old)

	details, err := ScanPath(tmp)
	if err != nil {
		t.Fatal(err)
	}

	var hits int
	for _, d := range details {
		if d.Sub == SubArchiveExtracted && d.Item.Name == "Project.zip" {
			hits++
		}
	}
	if hits != 1 {
		t.Fatalf("expected exactly 1 archive-extracted hit, got %d (details=%+v)", hits, details)
	}
}

// TestScanPath_DetectsProjectFolder: a directory containing node_modules.
func TestScanPath_DetectsProjectFolder(t *testing.T) {
	tmp := t.TempDir()
	mod := time.Now().AddDate(0, 0, -5)

	projDir := filepath.Join(tmp, "MyProject")
	makeDir(t, projDir, mod)
	makeFile(t, filepath.Join(projDir, "node_modules", "x"), 1024, mod)

	details, err := ScanPath(tmp)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, d := range details {
		if d.Sub == SubProjectFolder && d.Item.Name == "MyProject" {
			found = true
			if d.Item.Risk != item.RiskAskBefore {
				t.Errorf("project folder should be RiskAskBefore, got %v", d.Item.Risk)
			}
		}
	}
	if !found {
		t.Fatalf("expected project-folder detection; got %+v", details)
	}
}

// TestScanPath_OldDBDump: a 40-day-old .sql file is flagged as db-dump.
func TestScanPath_OldDBDump(t *testing.T) {
	tmp := t.TempDir()
	old := time.Now().AddDate(0, 0, -40)
	makeFile(t, filepath.Join(tmp, "backup.sql"), 1024, old)

	details, err := ScanPath(tmp)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range details {
		if d.Sub == SubDBDump && d.Item.Name == "backup.sql" {
			return
		}
	}
	t.Fatalf("expected db-dump detection; got %+v", details)
}

// TestScanPath_RecentDBDumpNotFlagged: a 5-day-old .sql is too fresh; should
// fall through to large-other only if >100 MB. With 1 KB it's nothing.
func TestScanPath_RecentDBDumpNotFlagged(t *testing.T) {
	tmp := t.TempDir()
	fresh := time.Now().AddDate(0, 0, -5)
	makeFile(t, filepath.Join(tmp, "fresh.sql"), 1024, fresh)

	details, err := ScanPath(tmp)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range details {
		if d.Sub == SubDBDump {
			t.Fatalf("fresh sql should NOT be flagged as db-dump; got %+v", d)
		}
	}
}

// TestScanPath_LargeOther: a 200 MB unclassified file shows up in large-other.
func TestScanPath_LargeOther(t *testing.T) {
	tmp := t.TempDir()
	mod := time.Now().AddDate(0, 0, -1)
	makeFile(t, filepath.Join(tmp, "mystery.bin"), 200*1024*1024, mod)

	details, err := ScanPath(tmp)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range details {
		if d.Sub == SubLargeOther && d.Item.Name == "mystery.bin" {
			return
		}
	}
	t.Fatalf("expected large-other for 200MB file; got %+v", details)
}

// TestScanPath_MissingFolder: scanning a non-existent path returns no items, no error.
func TestScanPath_MissingFolder(t *testing.T) {
	details, err := ScanPath("/tmp/this-path-must-not-exist-chipawa-12345")
	if err != nil {
		t.Fatalf("missing folder should not error, got %v", err)
	}
	if len(details) != 0 {
		t.Fatalf("expected zero details, got %d", len(details))
	}
}

// TestAsItems_ExcludesLargeOther: the cleaner feed must skip large-other entries.
func TestAsItems_ExcludesLargeOther(t *testing.T) {
	details := []Detail{
		{Sub: SubInstaller, Item: item.Item{Name: "a.dmg"}},
		{Sub: SubLargeOther, Item: item.Item{Name: "mystery.bin"}},
		{Sub: SubProjectFolder, Item: item.Item{Name: "proj"}},
	}
	items := AsItems(details)
	if len(items) != 2 {
		t.Fatalf("expected 2 items (large-other excluded), got %d", len(items))
	}
	for _, it := range items {
		if it.Name == "mystery.bin" {
			t.Errorf("large-other should have been excluded")
		}
	}
}
