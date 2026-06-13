// Package cleaner removes items previously detected by scanners.
//
// Design:
//   - A Remover knows how to delete one Item (path-based, docker-prune, etc.).
//   - A Plan groups items + their resolved Remover.
//   - Mode controls confirmation flow: DryRun, Interactive, Yes.
//   - Prompter abstracts the user interaction so we can test it.
//
// Safety rules enforced here:
//   - Never delete a path outside a known-safe prefix (home dir or /var/folders).
//   - DryRun never touches the filesystem.
//   - Docker volumes are NEVER pruned (only `system prune -f`, no --volumes).
//   - Empty/unknown paths are skipped, not errored, to avoid panics on malformed Items.
package cleaner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sistematlan/mistah/internal/item"
)

// Mode controls how the cleaner asks for confirmation.
type Mode int

const (
	// DryRun reports what would be removed but never touches disk.
	DryRun Mode = iota
	// Interactive asks the user before each item.
	Interactive
	// Yes assumes confirmation for every item (CI / scripting).
	Yes
)

// Decision is what the user (or auto mode) chose for an item.
type Decision int

const (
	DecisionYes Decision = iota
	DecisionNo
	DecisionView
	DecisionQuit
)

// Prompter renders an item and reads the user's decision.
// It is an interface so tests can inject a deterministic answer stream.
type Prompter interface {
	Ask(it item.Item) Decision
	Show(msg string)
}

// Remover deletes the resource an Item points to.
// Implementations must be idempotent and safe to call after a missing target.
type Remover interface {
	// Describe returns a one-line human description used in dry-run output.
	Describe(it item.Item) string
	// Remove performs the actual deletion. It must not call os.Exit.
	Remove(it item.Item) error
}

// Resolver picks a Remover for an item. The default resolver maps:
//   - Docker reclaimable → DockerPruneRemover
//   - everything else with a non-empty Path → PathRemover
type Resolver func(it item.Item) (Remover, error)

// DefaultResolver is the standard mapping used by `mistah clean`.
func DefaultResolver(it item.Item) (Remover, error) {
	if it.Tool == "docker" && it.Path == "" {
		return DockerPruneRemover{}, nil
	}
	if it.Path == "" {
		return nil, fmt.Errorf("item %q has no Path and no specialized remover", it.Name)
	}
	return PathRemover{}, nil
}

// Result reports what happened to one item after Run.
type Result struct {
	Item    item.Item
	Skipped bool   // user said no, or DryRun
	Error   error  // non-nil if removal failed
	Reason  string // why it was skipped
}

// Bytes returns the size that was (or would be) freed by this result.
func (r Result) Bytes() int64 {
	if r.Skipped || r.Error != nil {
		return 0
	}
	return r.Item.Bytes
}

// Plan groups the items the user chose to clean.
// A Plan is created from a list of detected items + a Mode + a Prompter.
type Plan struct {
	Items    []item.Item
	Mode     Mode
	Prompter Prompter
	Resolver Resolver
	Out      io.Writer // where dry-run lines and progress go
}

// New builds a plan with sensible defaults. Pass a nil Prompter for non-interactive use.
func New(items []item.Item, mode Mode, p Prompter, out io.Writer) *Plan {
	if out == nil {
		out = os.Stdout
	}
	return &Plan{
		Items:    items,
		Mode:     mode,
		Prompter: p,
		Resolver: DefaultResolver,
		Out:      out,
	}
}

// Run iterates over the plan respecting Mode and Prompter.
// It returns one Result per item, in the same order.
func (p *Plan) Run() []Result {
	results := make([]Result, 0, len(p.Items))
	for _, it := range p.Items {
		// Resolve remover before asking — if the item is malformed we skip with reason.
		remover, err := p.Resolver(it)
		if err != nil {
			results = append(results, Result{Item: it, Skipped: true, Reason: err.Error()})
			continue
		}

		decision := p.decide(it, remover)
		if decision == DecisionQuit {
			results = append(results, Result{Item: it, Skipped: true, Reason: "user quit"})
			break
		}
		if decision == DecisionNo {
			results = append(results, Result{Item: it, Skipped: true, Reason: "user declined"})
			continue
		}

		if p.Mode == DryRun {
			fmt.Fprintf(p.Out, "[dry-run] would remove: %s\n", remover.Describe(it))
			results = append(results, Result{Item: it, Skipped: true, Reason: "dry-run"})
			continue
		}

		fmt.Fprintf(p.Out, "removing %s ... ", it.Name)
		if err := remover.Remove(it); err != nil {
			fmt.Fprintf(p.Out, "FAILED: %v\n", err)
			results = append(results, Result{Item: it, Error: err})
			continue
		}
		fmt.Fprintln(p.Out, "ok")
		results = append(results, Result{Item: it})
	}
	return results
}

// decide returns the user decision for an item according to Mode.
//
// In Yes mode we always return DecisionYes.
// In DryRun mode we always return DecisionYes too — the caller will skip
// the actual deletion but still report what *would* happen.
// In Interactive mode we go through the Prompter and handle the "view" loop.
func (p *Plan) decide(it item.Item, remover Remover) Decision {
	if p.Mode == Yes || p.Mode == DryRun {
		return DecisionYes
	}
	if p.Prompter == nil {
		return DecisionNo // safest fallback for misuse
	}
	for {
		d := p.Prompter.Ask(it)
		if d != DecisionView {
			return d
		}
		// Show extra context (path listing) and ask again.
		p.Prompter.Show(viewItem(it, remover))
	}
}

// viewItem builds the multi-line preview shown when the user types "v".
// For PathRemover we list the immediate children. For Docker we run `docker system df`.
func viewItem(it item.Item, remover Remover) string {
	switch r := remover.(type) {
	case PathRemover:
		return r.previewPath(it.Path)
	case DockerPruneRemover:
		out, _ := exec.Command("docker", "system", "df").Output()
		return string(out)
	default:
		return fmt.Sprintf("(no preview available for %T)", remover)
	}
}

// SafeRoots are the only filesystem prefixes a PathRemover will touch.
// Anything outside this set is rejected with ErrUnsafePath.
//
// /var/folders and /tmp are included so tests using TempDir work transparently;
// they are also legitimate caches in macOS.
var SafeRoots = func() []string {
	roots := []string{"/var/folders", "/tmp", "/private/var/folders", "/private/tmp"}
	if home, err := os.UserHomeDir(); err == nil {
		roots = append(roots, home)
	}
	return roots
}()

// ErrUnsafePath is returned when a PathRemover is asked to delete outside SafeRoots.
var ErrUnsafePath = errors.New("path is outside safe roots; refusing to delete")

// PathRemover deletes a directory or file with rm -rf semantics, but only if
// the path is rooted in SafeRoots. Use this for cache directories.
type PathRemover struct{}

func (PathRemover) Describe(it item.Item) string {
	return fmt.Sprintf("%s (%s)", it.Path, it.Name)
}

func (PathRemover) Remove(it item.Item) error {
	if it.Path == "" {
		return errors.New("empty path")
	}
	abs, err := filepath.Abs(it.Path)
	if err != nil {
		return err
	}
	if !isUnderSafeRoot(abs) {
		return fmt.Errorf("%w: %s", ErrUnsafePath, abs)
	}
	return os.RemoveAll(abs)
}

// previewPath returns up to N child entries with sizes for the "view" UI.
func (PathRemover) previewPath(path string) string {
	const maxEntries = 15
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Sprintf("(cannot read %s: %v)", path, err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Contents of %s:\n", path)
	limit := len(entries)
	if limit > maxEntries {
		limit = maxEntries
	}
	for i := 0; i < limit; i++ {
		fmt.Fprintf(&b, "  %s\n", entries[i].Name())
	}
	if len(entries) > maxEntries {
		fmt.Fprintf(&b, "  ... and %d more\n", len(entries)-maxEntries)
	}
	return b.String()
}

// isUnderSafeRoot returns true iff abs starts with one of the SafeRoots
// at a path-component boundary (so /home/foo doesn't match /homer).
func isUnderSafeRoot(abs string) bool {
	for _, root := range SafeRoots {
		if abs == root {
			return true
		}
		if strings.HasPrefix(abs, root+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

// DockerPruneRemover invokes `docker system prune -f` to free reclaimable space.
// It does NOT pass --volumes; user data on volumes stays untouched.
type DockerPruneRemover struct{}

func (DockerPruneRemover) Describe(it item.Item) string {
	return "docker system prune -f (images, build cache, stopped containers; volumes preserved)"
}

func (DockerPruneRemover) Remove(it item.Item) error {
	cmd := exec.Command("docker", "system", "prune", "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker prune failed: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Summary aggregates Run results for the final report.
type Summary struct {
	Removed     int
	Skipped     int
	Failed      int
	BytesFreed  int64
	BytesPlanned int64
}

// Summarize folds a slice of results into a Summary.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		s.BytesPlanned += r.Item.Bytes
		switch {
		case r.Error != nil:
			s.Failed++
		case r.Skipped:
			s.Skipped++
		default:
			s.Removed++
			s.BytesFreed += r.Item.Bytes
		}
	}
	return s
}
