package projects

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sistematlan/mistah/internal/disk"
)

type Project struct {
	Name       string
	Path       string
	Bytes      int64
	HasGit     bool
	Remote     string
	LastCommit time.Time
	NoCommits  bool
}

func (p Project) HasGitStr() string {
	if p.HasGit {
		return "sí"
	}
	return "no"
}

func (p Project) LastCommitLabel() string {
	if !p.HasGit || p.NoCommits {
		return "—"
	}
	days := int(time.Since(p.LastCommit).Hours() / 24)
	if days == 0 {
		return "hoy"
	}
	if days < 365 {
		return p.LastCommit.Format("2006-01-02")
	}
	return p.LastCommit.Format("2006-01-02") + " ⚠"
}

func (p Project) ShortRemote() string {
	if p.Remote == "" {
		return "—"
	}
	// Strip protocol prefix for display
	r := p.Remote
	r = strings.TrimPrefix(r, "https://")
	r = strings.TrimPrefix(r, "git@")
	r = strings.ReplaceAll(r, ":", "/")
	if len(r) > 30 {
		return r[:27] + "..."
	}
	return r
}

func Scan(root string) ([]Project, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var list []Project
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(root, e.Name())
		bytes, _ := disk.DirSize(path)

		p := Project{
			Name:  e.Name(),
			Path:  path,
			Bytes: bytes,
		}

		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			p.HasGit = true
			p.Remote = gitRemote(path)
			p.LastCommit, p.NoCommits = lastCommit(path)
		}

		list = append(list, p)
	}
	return list, nil
}

func gitRemote(path string) string {
	out, err := exec.Command("git", "-C", path, "remote", "get-url", "origin").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func lastCommit(path string) (time.Time, bool) {
	out, err := exec.Command("git", "-C", path, "log", "-1", "--format=%aI").Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return time.Time{}, true
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
	if err != nil {
		return time.Time{}, true
	}
	return t, false
}
