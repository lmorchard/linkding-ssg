package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env")

	content := `
# Comment line
LINKDING_URL=https://my-linkding.test
export LINKDING_TOKEN="secret-token-123"
OTHER_VAL='quoted single'
PRE_EXISTING=new-val
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Clean environment variables that could leak from runner
	origURL, hasURL := os.LookupEnv("LINKDING_URL")
	origToken, hasToken := os.LookupEnv("LINKDING_TOKEN")
	os.Unsetenv("LINKDING_URL")
	os.Unsetenv("LINKDING_TOKEN")
	defer func() {
		if hasURL {
			os.Setenv("LINKDING_URL", origURL)
		} else {
			os.Unsetenv("LINKDING_URL")
		}
		if hasToken {
			os.Setenv("LINKDING_TOKEN", origToken)
		} else {
			os.Unsetenv("LINKDING_TOKEN")
		}
	}()

	// Set pre-existing to test precedence
	os.Setenv("PRE_EXISTING", "original-val")
	defer os.Unsetenv("PRE_EXISTING")
	defer os.Unsetenv("OTHER_VAL")

	if err := LoadEnv(envPath); err != nil {
		t.Fatalf("LoadEnv failed: %v", err)
	}

	if got := os.Getenv("LINKDING_URL"); got != "https://my-linkding.test" {
		t.Errorf("expected https://my-linkding.test, got %q", got)
	}
	if got := os.Getenv("LINKDING_TOKEN"); got != "secret-token-123" {
		t.Errorf("expected secret-token-123, got %q", got)
	}
	if got := os.Getenv("OTHER_VAL"); got != "quoted single" {
		t.Errorf("expected 'quoted single', got %q", got)
	}
	if got := os.Getenv("PRE_EXISTING"); got != "original-val" {
		t.Errorf("expected pre-existing env to win, got %q", got)
	}
}

func TestLoadEnvMissingFile(t *testing.T) {
	err := LoadEnv("non-existent-file.env")
	if err != nil {
		t.Fatalf("expected missing file to be ignored, got err: %v", err)
	}
}
