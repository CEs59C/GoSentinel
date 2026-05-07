package config

import (
	"os"
	"testing"
)

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("POST_IN", "test@in.com")
	os.Setenv("POST_TO", "test@to.com")
	os.Setenv("PASSWORD", "123")
	os.Setenv("MODE", "console")

	defer os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.From != "test@in.com" {
		t.Error("POST_IN not loaded correctly")
	}

	if cfg.Mode != ModeConsole {
		t.Error("mode not parsed correctly")
	}
}
