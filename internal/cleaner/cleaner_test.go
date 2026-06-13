package cleaner

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/sistematlan/mistah/internal/item"
)

// scriptedPrompter returns a fixed sequence of decisions, one per Ask call.
// Tests fail loudly if Ask is called more times than answers were provided.
type scriptedPrompter struct {
	answers []Decision
	calls   int
	shown   []string
}

func (p *scriptedPrompter) Ask(_ item.Item) Decision {
	if p.calls >= len(p.answers) {
		panic("scriptedPrompter: Ask called more times than answers")
	}
	d := p.answers[p.calls]
	p.calls++
	return d
}
func (p *scriptedPrompter) Show(msg string) { p.shown = append(p.shown, msg) }

// makeFile creates a temp directory with one byte of content so DirSize > 0.
func makeFile(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestDryRun_DoesNotDelete: dry-run never touches disk, even with Yes answers.
func TestDryRun_DoesNotDelete(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "cache")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	makeFile(t, target, "a.txt")

	it := item.Item{Name: "test", Path: target, Bytes: 1, Risk: item.RiskSafe}
	plan := New([]item.Item{it}, DryRun, &scriptedPrompter{answers: []Decision{DecisionYes}}, &bytes.Buffer{})
	results := plan.Run()

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("dry-run deleted target: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Fatalf("expected 1 skipped result, got %+v", results)
	}
}

// TestYesMode_DeletesWithoutPrompt: Yes mode removes every item, no prompter calls.
func TestYesMode_DeletesWithoutPrompt(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(tmp, "b")
	for _, p := range []string{a, b} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
		makeFile(t, p, "f")
	}
	items := []item.Item{
		{Name: "a", Path: a, Bytes: 10, Risk: item.RiskSafe},
		{Name: "b", Path: b, Bytes: 20, Risk: item.RiskSafe},
	}
	// nil prompter is fine in Yes mode — it shouldn't be called.
	plan := New(items, Yes, nil, &bytes.Buffer{})
	results := plan.Run()

	if _, err := os.Stat(a); !os.IsNotExist(err) {
		t.Errorf("a should be gone, got err=%v", err)
	}
	if _, err := os.Stat(b); !os.IsNotExist(err) {
		t.Errorf("b should be gone, got err=%v", err)
	}
	s := Summarize(results)
	if s.Removed != 2 || s.BytesFreed != 30 {
		t.Fatalf("Summary = %+v", s)
	}
}

// TestInteractive_NoSkips: a "no" answer skips that item, others continue.
func TestInteractive_NoSkips(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(tmp, "b")
	for _, p := range []string{a, b} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	items := []item.Item{
		{Name: "a", Path: a, Bytes: 1, Risk: item.RiskSafe},
		{Name: "b", Path: b, Bytes: 1, Risk: item.RiskSafe},
	}
	pr := &scriptedPrompter{answers: []Decision{DecisionNo, DecisionYes}}
	plan := New(items, Interactive, pr, &bytes.Buffer{})
	results := plan.Run()

	if _, err := os.Stat(a); err != nil {
		t.Errorf("a should still exist (declined)")
	}
	if _, err := os.Stat(b); !os.IsNotExist(err) {
		t.Errorf("b should be gone (accepted)")
	}
	if results[0].Skipped != true || results[1].Skipped != false {
		t.Fatalf("results = %+v", results)
	}
}

// TestInteractive_QuitStopsPlan: a quit answer halts processing of remaining items.
func TestInteractive_QuitStopsPlan(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(tmp, "b")
	for _, p := range []string{a, b} {
		_ = os.MkdirAll(p, 0o755)
	}
	items := []item.Item{
		{Name: "a", Path: a, Bytes: 1, Risk: item.RiskSafe},
		{Name: "b", Path: b, Bytes: 1, Risk: item.RiskSafe},
	}
	pr := &scriptedPrompter{answers: []Decision{DecisionQuit}}
	plan := New(items, Interactive, pr, &bytes.Buffer{})
	results := plan.Run()

	if len(results) != 1 {
		t.Fatalf("expected processing to stop after quit, got %d results", len(results))
	}
	if pr.calls != 1 {
		t.Fatalf("prompter called %d times after quit", pr.calls)
	}
}

// TestInteractive_ViewLoopThenYes: view triggers Show, then next Ask is consulted.
func TestInteractive_ViewLoopThenYes(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	if err := os.MkdirAll(a, 0o755); err != nil {
		t.Fatal(err)
	}
	makeFile(t, a, "child.txt")

	items := []item.Item{{Name: "a", Path: a, Bytes: 1, Risk: item.RiskSafe}}
	pr := &scriptedPrompter{answers: []Decision{DecisionView, DecisionYes}}
	plan := New(items, Interactive, pr, &bytes.Buffer{})
	plan.Run()

	if _, err := os.Stat(a); !os.IsNotExist(err) {
		t.Errorf("a should be removed after view+yes")
	}
	if len(pr.shown) != 1 {
		t.Errorf("Show called %d times, want 1", len(pr.shown))
	}
}

// TestPathRemover_RejectsUnsafePath: deletion outside SafeRoots must error.
func TestPathRemover_RejectsUnsafePath(t *testing.T) {
	r := PathRemover{}
	err := r.Remove(item.Item{Path: "/etc/hosts"})
	if err == nil {
		t.Fatal("expected error for unsafe path")
	}
}

// TestPathRemover_EmptyPath: empty Path is a hard error, never a panic.
func TestPathRemover_EmptyPath(t *testing.T) {
	if err := (PathRemover{}).Remove(item.Item{}); err == nil {
		t.Fatal("expected error for empty path")
	}
}

// TestDefaultResolver_PicksDocker: docker tool with no Path → DockerPruneRemover.
func TestDefaultResolver_PicksDocker(t *testing.T) {
	r, err := DefaultResolver(item.Item{Tool: "docker", Path: ""})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(DockerPruneRemover); !ok {
		t.Fatalf("expected DockerPruneRemover, got %T", r)
	}
}

// TestDefaultResolver_PicksPath: any item with a Path → PathRemover.
func TestDefaultResolver_PicksPath(t *testing.T) {
	r, err := DefaultResolver(item.Item{Tool: "npm", Path: "/x"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.(PathRemover); !ok {
		t.Fatalf("expected PathRemover, got %T", r)
	}
}

// TestDefaultResolver_RejectsMalformed: empty Path + non-docker tool → error.
func TestDefaultResolver_RejectsMalformed(t *testing.T) {
	if _, err := DefaultResolver(item.Item{Tool: "npm", Path: ""}); err == nil {
		t.Fatal("expected resolver to reject empty-path non-docker item")
	}
}
