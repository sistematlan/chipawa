// Package spinner provides a minimal terminal spinner for tasks where the
// user has nothing to look at while we shell out to du or docker.
//
// Design:
//   - One goroutine flips frames every Interval. No heap of dependencies.
//   - Detects non-TTY output and degrades to a one-shot "[…]" line so
//     pipes and redirected output stay readable.
//   - The Frames are unicode braille dots to match common Linux/macOS look.
//
// Usage:
//
//	sp := spinner.New("Scanning your disk")
//	sp.Start(os.Stdout)
//	defer sp.Stop()
//	... long work ...
//
// The same spinner can be reused: Stop() then Start() with a new label.
package spinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Frames is the rotating glyph set. 10-frame braille dots animate smoothly
// at 80ms per frame.
var Frames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// DefaultInterval is the per-frame delay; chosen to feel alive but not jittery.
const DefaultInterval = 80 * time.Millisecond

// Spinner shows a rotating glyph next to a label until Stop() is called.
type Spinner struct {
	Label    string
	Interval time.Duration

	mu      sync.Mutex
	out     io.Writer
	stop    chan struct{}
	wg      sync.WaitGroup
	running bool
	tty     bool
}

// New creates a spinner with the given label. Call Start to render.
func New(label string) *Spinner {
	return &Spinner{Label: label, Interval: DefaultInterval}
}

// Start begins the animation. Safe to call once per Spinner instance per cycle.
// If out is not a terminal, Start prints a single static line and returns.
func (s *Spinner) Start(out io.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return
	}
	s.out = out
	s.tty = isTerminal(out)
	s.stop = make(chan struct{})
	s.running = true

	if !s.tty {
		// In a pipe / redirected output, spinning would dump garbage. Print
		// a single status line and skip the goroutine.
		fmt.Fprintf(out, "[…] %s\n", s.Label)
		return
	}

	s.wg.Add(1)
	go s.loop()
}

// Stop halts the animation and clears the line. Safe to call when not running.
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	close(s.stop)
	tty := s.tty
	out := s.out
	s.mu.Unlock()

	s.wg.Wait()

	if tty && out != nil {
		// Erase line then move cursor to column 1 so the next print starts clean.
		fmt.Fprint(out, "\r\033[2K")
	}
}

// SetLabel changes the spinner text mid-animation. Useful when scanning
// multiple sections sequentially.
func (s *Spinner) SetLabel(label string) {
	s.mu.Lock()
	s.Label = label
	s.mu.Unlock()
}

// loop is the animation goroutine.
func (s *Spinner) loop() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	frame := 0
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.mu.Lock()
			label := s.Label
			out := s.out
			s.mu.Unlock()
			// \r returns to column 1; \033[2K clears the line so we don't
			// leave debris when the label shrinks.
			fmt.Fprintf(out, "\r\033[2K%s %s", Frames[frame%len(Frames)], label)
			frame++
		}
	}
}

// isTerminal reports whether w is a TTY. We avoid pulling x/term in just
// for this — a Stat() on os.Stdout is enough on macOS/Linux.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
