package apps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sistematlan/mistah/internal/disk"
)

type App struct {
	Name      string
	Path      string
	Bytes     int64
	LastUsed  time.Time
	NeverUsed bool
	DaysSinceUse int
}

func (a App) LastUsedLabel() string {
	if a.NeverUsed {
		return "nunca"
	}
	if a.DaysSinceUse == 0 {
		return "hoy"
	}
	if a.DaysSinceUse < 30 {
		return fmt.Sprintf("hace %d días", a.DaysSinceUse)
	}
	return a.LastUsed.Format("2006-01-02")
}

func List() ([]App, error) {
	dirs := []string{"/Applications"}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(home, "Applications"))
	}

	var apps []App
	seen := map[string]bool{}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".app") {
				continue
			}
			name := strings.TrimSuffix(e.Name(), ".app")
			if seen[name] {
				continue
			}
			seen[name] = true

			path := filepath.Join(dir, e.Name())
			bytes, _ := disk.DirSize(path)
			lastUsed, neverUsed := lastUsedDate(path)
			days := -1
			if !neverUsed {
				days = int(time.Since(lastUsed).Hours() / 24)
			}

			apps = append(apps, App{
				Name:         name,
				Path:         path,
				Bytes:        bytes,
				LastUsed:     lastUsed,
				NeverUsed:    neverUsed,
				DaysSinceUse: days,
			})
		}
	}

	return apps, nil
}

func lastUsedDate(path string) (time.Time, bool) {
	out, err := exec.Command("mdls", "-name", "kMDItemLastUsedDate", "-raw", path).Output()
	if err != nil {
		return time.Time{}, true
	}
	raw := strings.TrimSpace(string(out))
	if raw == "(null)" || raw == "" {
		return time.Time{}, true
	}
	t, err := time.Parse("2006-01-02 15:04:05 +0000", raw)
	if err != nil {
		return time.Time{}, true
	}
	return t, false
}
