package ui

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	sharedCompat "image-studio/shared/compat"
)

func TestCopyImageFileCopiesToExplicitPath(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.png")
	if err := os.WriteFile(src, []byte("image"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}
	dst := filepath.Join(dir, "nested", "copy.png")
	saved, err := copyImageFile(src, dst)
	if err != nil {
		t.Fatalf("copyImageFile: %v", err)
	}
	if saved != dst {
		t.Fatalf("saved=%q want %q", saved, dst)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copied: %v", err)
	}
	if string(data) != "image" {
		t.Fatalf("copied data=%q", data)
	}
}

func TestCopyImageFileDirectoryTargetKeepsSourceName(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "source.webp")
	if err := os.WriteFile(src, []byte("image"), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}
	targetDir := filepath.Join(dir, "target")
	if err := os.Mkdir(targetDir, 0o700); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	saved, err := copyImageFile(src, targetDir)
	if err != nil {
		t.Fatalf("copyImageFile: %v", err)
	}
	want := filepath.Join(targetDir, "source.webp")
	if saved != want {
		t.Fatalf("saved=%q want %q", saved, want)
	}
}

func TestMatchHistoryQueryMatchesPromptAndPath(t *testing.T) {
	item := sharedCompat.HistoryItem{
		Prompt:        "生成一张雪山海报",
		RevisedPrompt: "cinematic snow mountain poster",
		SavedPath:     "/tmp/snow.png",
		Size:          "1024x1024",
		Quality:       "high",
	}
	if !matchHistoryQuery(item, "雪山") {
		t.Fatalf("expected prompt match")
	}
	if !matchHistoryQuery(item, "snow.png") {
		t.Fatalf("expected path match")
	}
	if matchHistoryQuery(item, "desert") {
		t.Fatalf("unexpected query match")
	}
}

func TestTodayHistoryCountUsesLocalDayBoundary(t *testing.T) {
	now := time.Date(2026, time.May, 31, 15, 4, 0, 0, time.Local)
	items := []sharedCompat.HistoryItem{
		{ID: "a", CreatedAt: now.Add(-2 * time.Hour).UnixMilli()},
		{ID: "b", CreatedAt: now.Add(-26 * time.Hour).UnixMilli()},
	}
	if got := todayHistoryCount(items, now); got != 1 {
		t.Fatalf("todayHistoryCount=%d want 1", got)
	}
}
