package orphans

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWhatsappMedia_Empty verifies that with no media folder, no item is reported.
func TestWhatsappMedia_Empty(t *testing.T) {
	tmp := t.TempDir()
	if _, ok := whatsappMedia(tmp); ok {
		t.Fatalf("expected no item when Media folder is missing")
	}
}

// TestWhatsappMedia_Detected verifies that a populated media folder is reported.
func TestWhatsappMedia_Detected(t *testing.T) {
	tmp := t.TempDir()
	mediaDir := filepath.Join(tmp,
		"Library/Group Containers/group.net.whatsapp.WhatsApp.shared/Message/Media")
	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Write something measurable so DirSize > 0.
	if err := os.WriteFile(filepath.Join(mediaDir, "fake.jpg"), make([]byte, 4096), 0o644); err != nil {
		t.Fatal(err)
	}

	it, ok := whatsappMedia(tmp)
	if !ok {
		t.Fatalf("expected an item to be detected")
	}
	if it.Path != mediaDir {
		t.Errorf("Path = %q; want %q", it.Path, mediaDir)
	}
	if it.Bytes <= 0 {
		t.Errorf("Bytes = %d; want > 0", it.Bytes)
	}
	if it.Tool != "whatsapp" {
		t.Errorf("Tool = %q; want whatsapp", it.Tool)
	}
}

// TestDockerLeftover_Empty: no container dir → no item.
func TestDockerLeftover_Empty(t *testing.T) {
	tmp := t.TempDir()
	if _, ok := dockerLeftover(tmp); ok {
		t.Fatalf("expected no item when container dir is missing")
	}
}
