package kernel

import (
	"path/filepath"
	"testing"

	"github.com/yuanhua/image-gptcodex/pkg/client"
)

func TestParseSourcePaths(t *testing.T) {
	got := ParseSourcePaths(" /tmp/a.png\n'/tmp/b.jpg',\"/tmp/a.png\" ")
	want := []string{"/tmp/a.png", "/tmp/b.jpg"}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%q want %q", i, got[i], want[i])
		}
	}
}

func TestNormalizeConfigDefaults(t *testing.T) {
	cfg := normalizeConfig(Config{
		Prompt:    "  hello  ",
		Mode:      client.Mode("unknown"),
		OutputDir: filepath.Join("tmp", "out"),
	})
	if cfg.Prompt != "hello" {
		t.Fatalf("prompt=%q", cfg.Prompt)
	}
	if cfg.Mode != client.ModeGenerate {
		t.Fatalf("mode=%q", cfg.Mode)
	}
	if cfg.APIMode != client.APIModeResponses {
		t.Fatalf("api mode=%q", cfg.APIMode)
	}
	if cfg.TextModelID == "" || cfg.ImageModelID == "" || cfg.OutputFormat == "" {
		t.Fatalf("missing defaults: %#v", cfg)
	}
}

func TestBuildImageNameMapsJPEGExtension(t *testing.T) {
	got := buildImageName(client.ModeEdit, "A cat wearing sunglasses", "20260531-120000", "jpeg")
	want := "image-edit-a-cat-wearing-sunglasses-20260531-120000.jpg"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
