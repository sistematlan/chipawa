// Package orphans detects user data that belongs to apps no longer installed,
// or media caches that grew without bound. These items are NOT regular caches:
// they may hold user-visible content, so detectors mark them RiskAskBefore.
package orphans

import (
	"os"
	"path/filepath"

	"github.com/sistematlan/mistah/internal/disk"
	"github.com/sistematlan/mistah/internal/item"
)

// Scan inspects the system for orphaned data and returns the items found.
func Scan() ([]item.Item, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var items []item.Item

	if it, ok := dockerLeftover(home); ok {
		items = append(items, it)
	}
	if it, ok := whatsappMedia(home); ok {
		items = append(items, it)
	}

	return items, nil
}

// dockerLeftover finds the Docker Desktop container directory left behind
// when the app was uninstalled but its 30+ GB of VM data remains on disk.
//
// Heuristic: if /Applications/Docker.app does NOT exist but
// ~/Library/Containers/com.docker.docker does, it's leftover.
func dockerLeftover(home string) (item.Item, bool) {
	containerPath := filepath.Join(home, "Library/Containers/com.docker.docker")
	if _, err := os.Stat(containerPath); err != nil {
		return item.Item{}, false
	}
	if _, err := os.Stat("/Applications/Docker.app"); err == nil {
		return item.Item{}, false // Docker is installed; not orphan
	}
	bytes, _ := disk.DirSize(containerPath)
	if bytes <= 0 {
		return item.Item{}, false
	}
	return item.Item{
		Name:      "Docker Desktop leftover",
		NameKey:   "orphans.docker-leftover.name",
		Tool:      "docker",
		Path:      containerPath,
		Bytes:     bytes,
		Category:  item.CategoryOrphan,
		Risk:      item.RiskAskBefore,
		Detail:    "Docker.app is uninstalled but its container data remains",
		DetailKey: "orphans.docker-leftover.detail",
	}, true
}

// whatsappMedia finds media cached by WhatsApp Desktop. Removing it does
// NOT delete chats — only photos, videos and audio that re-download on demand.
//
// We only target the Media subfolder so ChatStorage.sqlite stays intact.
func whatsappMedia(home string) (item.Item, bool) {
	mediaPath := filepath.Join(home,
		"Library/Group Containers/group.net.whatsapp.WhatsApp.shared/Message/Media")
	bytes, _ := disk.DirSize(mediaPath)
	if bytes <= 0 {
		return item.Item{}, false
	}
	return item.Item{
		Name:      "WhatsApp media",
		NameKey:   "orphans.whatsapp-media.name",
		Tool:      "whatsapp",
		Path:      mediaPath,
		Bytes:     bytes,
		Category:  item.CategoryOrphan,
		Risk:      item.RiskAskBefore,
		Detail:    "downloaded photos/videos/audio (chats not affected)",
		DetailKey: "orphans.whatsapp-media.detail",
	}, true
}
