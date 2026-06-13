package cleaner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/i18n"
	"github.com/sistematlan/mistah/internal/item"
)

// SimpleMode toggles human-friendly phrasing in prompts. The cmd package
// sets it before each Run; tests can leave it as the zero value (advanced).
var SimpleMode bool

// TerminalPrompter is the default Prompter used by `mistah clean`.
// It reads a single line from stdin per Ask call and accepts:
//
//	s / si / y / yes  → DecisionYes
//	n / no            → DecisionNo  (default on empty)
//	v / ver / view    → DecisionView (mistah shows path contents and asks again)
//	q / quit          → DecisionQuit (stop the whole plan)
//
// For RiskDangerous items the prompter requires the user to TYPE the item name
// to confirm. This prevents accidental fat-finger destruction.
//
// IMPORTANT: a single bufio.Scanner is reused across calls. Allocating a new
// Scanner per Ask would discard buffered bytes between prompts and silently
// drop the user's input on the next round.
type TerminalPrompter struct {
	In      io.Reader
	Out     io.Writer
	scanner *bufio.Scanner
}

// NewTerminalPrompter returns a prompter wired to stdin/stdout.
func NewTerminalPrompter() *TerminalPrompter {
	return &TerminalPrompter{In: os.Stdin, Out: os.Stdout}
}

// reader returns a lazily initialised scanner shared across Ask calls.
func (p *TerminalPrompter) reader() *bufio.Scanner {
	if p.scanner == nil {
		p.scanner = bufio.NewScanner(p.In)
	}
	return p.scanner
}

func (p *TerminalPrompter) Show(msg string) {
	fmt.Fprintln(p.Out, msg)
}

// Ask prints the item and reads one line from In.
func (p *TerminalPrompter) Ask(it item.Item) Decision {
	p.printItem(it)
	if it.Risk == item.RiskDangerous {
		return p.askDangerous(it)
	}
	prompt := i18n.T("cleaner.prompt")
	return readDecision(p.reader(), p.Out, prompt)
}

// printItem renders the header block shown before each prompt.
func (p *TerminalPrompter) printItem(it item.Item) {
	fmt.Fprintln(p.Out)
	fmt.Fprintf(p.Out, "  %s — %s\n", it.HumanName(), disk.FormatBytes(it.Bytes))
	if it.Tool != "" {
		fmt.Fprintf(p.Out, "  %s : %s\n", i18n.T("ui.tool"), it.Tool)
	}
	if it.Path != "" && !SimpleMode {
		fmt.Fprintf(p.Out, "  %s : %s\n", i18n.T("ui.path"), it.Path)
	}
	if d := it.HumanDetail(SimpleMode); d != "" {
		fmt.Fprintf(p.Out, "  %s : %s\n", i18n.T("ui.note"), d)
	}
	fmt.Fprintf(p.Out, "  %s : %s\n", i18n.T("ui.risk"), it.HumanRisk())
}

// askDangerous forces the user to type the item name to confirm.
func (p *TerminalPrompter) askDangerous(it item.Item) Decision {
	fmt.Fprintf(p.Out, "  "+i18n.T("cleaner.prompt.dangerous"), it.HumanName())
	scanner := p.reader()
	if !scanner.Scan() {
		return DecisionNo
	}
	if strings.TrimSpace(scanner.Text()) == it.HumanName() {
		return DecisionYes
	}
	return DecisionNo
}

// readDecision parses one yes/no/view/quit answer.
// Empty input defaults to "no" (safest).
// The scanner is shared across calls; do NOT allocate a new one here.
func readDecision(scanner *bufio.Scanner, out io.Writer, prompt string) Decision {
	fmt.Fprint(out, "  "+prompt)
	if !scanner.Scan() {
		return DecisionNo
	}
	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "s", "si", "sí", "y", "yes":
		return DecisionYes
	case "v", "ver", "view":
		return DecisionView
	case "q", "quit", "salir":
		return DecisionQuit
	default:
		return DecisionNo
	}
}
