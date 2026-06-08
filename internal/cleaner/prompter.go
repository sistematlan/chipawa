package cleaner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sistematlan/chipawa/internal/disk"
	"github.com/sistematlan/chipawa/internal/item"
)

// TerminalPrompter is the default Prompter used by `chipawa clean`.
// It reads a single line from stdin per Ask call and accepts:
//
//	s / si / y / yes  → DecisionYes
//	n / no            → DecisionNo  (default on empty)
//	v / ver / view    → DecisionView (chipawa shows path contents and asks again)
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
	prompt := "[s/N/v=ver/q=salir] "
	return readDecision(p.reader(), p.Out, prompt)
}

// printItem renders the header block shown before each prompt.
func (p *TerminalPrompter) printItem(it item.Item) {
	fmt.Fprintln(p.Out)
	fmt.Fprintf(p.Out, "  %s — %s\n", it.Name, disk.FormatBytes(it.Bytes))
	if it.Tool != "" {
		fmt.Fprintf(p.Out, "  tool : %s\n", it.Tool)
	}
	if it.Path != "" {
		fmt.Fprintf(p.Out, "  path : %s\n", it.Path)
	}
	if it.Detail != "" {
		fmt.Fprintf(p.Out, "  note : %s\n", it.Detail)
	}
	fmt.Fprintf(p.Out, "  risk : %s\n", riskLabel(it.Risk))
}

// askDangerous forces the user to type the item name to confirm.
func (p *TerminalPrompter) askDangerous(it item.Item) Decision {
	fmt.Fprintf(p.Out, "  CONFIRMA escribiendo el nombre exacto (%q) o vacío para cancelar:\n  > ", it.Name)
	scanner := p.reader()
	if !scanner.Scan() {
		return DecisionNo
	}
	if strings.TrimSpace(scanner.Text()) == it.Name {
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

func riskLabel(r item.Risk) string {
	switch r {
	case item.RiskSafe:
		return "safe (cache regenerable)"
	case item.RiskAskBefore:
		return "ask-before (puede contener datos del usuario)"
	case item.RiskDangerous:
		return "dangerous (irreversible)"
	default:
		return string(r)
	}
}
