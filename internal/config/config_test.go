package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DownloadDir == "" {
		t.Error("DownloadDir should not be empty")
	}
	if !filepath.HasPrefix(cfg.DownloadDir, "/") && !filepath.HasPrefix(cfg.DownloadDir, "C:") {
		t.Errorf("DownloadDir should be absolute path, got: %s", cfg.DownloadDir)
	}
	if !strings.HasSuffix(cfg.DownloadDir, "Drift") {
		t.Errorf("DownloadDir should end with 'Drift', got: %s", cfg.DownloadDir)
	}

	if cfg.AcceptTimeout != 30*time.Second {
		t.Errorf("AcceptTimeout should be 30s, got: %v", cfg.AcceptTimeout)
	}

	if cfg.Identity != "" {
		t.Errorf("Identity should be empty string by default, got: %s", cfg.Identity)
	}
}

func TestDefaultPath(t *testing.T) {
	path := DefaultPath()

	if path == "" {
		t.Error("DefaultPath should not be empty")
	}
	if !strings.HasSuffix(path, "config.toml") {
		t.Errorf("DefaultPath should end with 'config.toml', got: %s", path)
	}
	if !filepath.IsAbs(path) {
		t.Errorf("DefaultPath should be absolute, got: %s", path)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.toml")

	if err != nil {
		t.Errorf("Load should not error on missing file, got: %v", err)
	}
	if cfg == nil {
		t.Error("Load should return config even on missing file")
	}

	// Should return defaults
	if cfg.AcceptTimeout != 30*time.Second {
		t.Errorf("Missing file should return default timeout, got: %v", cfg.AcceptTimeout)
	}
	if cfg.Identity != "" {
		t.Errorf("Missing file should return default identity, got: %s", cfg.Identity)
	}
}

func TestLoadValidFile(t *testing.T) {
	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "config.toml")

	content := `download_dir = "/tmp/MyDrift"
accept_timeout = "45s"
identity = "TestDevice"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)

	if err != nil {
		t.Errorf("Load should not error on valid file, got: %v", err)
	}
	if cfg == nil {
		t.Error("Load should return config")
	}

	if cfg.DownloadDir != "/tmp/MyDrift" {
		t.Errorf("DownloadDir should be '/tmp/MyDrift', got: %s", cfg.DownloadDir)
	}
	if cfg.AcceptTimeout != 45*time.Second {
		t.Errorf("AcceptTimeout should be 45s, got: %v", cfg.AcceptTimeout)
	}
	if cfg.Identity != "TestDevice" {
		t.Errorf("Identity should be 'TestDevice', got: %s", cfg.Identity)
	}
}

func TestLoadPartialFile(t *testing.T) {
	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "config.toml")

	content := `identity = "PartialDevice"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)

	if err != nil {
		t.Errorf("Load should not error on partial file, got: %v", err)
	}
	if cfg == nil {
		t.Error("Load should return config")
	}

	// Should have default download dir
	if cfg.DownloadDir == "" {
		t.Error("DownloadDir should have default value")
	}
	if !strings.HasSuffix(cfg.DownloadDir, "Drift") {
		t.Errorf("DownloadDir should end with 'Drift', got: %s", cfg.DownloadDir)
	}

	// Should have default timeout
	if cfg.AcceptTimeout != 30*time.Second {
		t.Errorf("AcceptTimeout should be default 30s, got: %v", cfg.AcceptTimeout)
	}

	// Should have value from file
	if cfg.Identity != "PartialDevice" {
		t.Errorf("Identity should be 'PartialDevice', got: %s", cfg.Identity)
	}
}

func TestLoadCorruptFile(t *testing.T) {
	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "config.toml")

	content := `this is not valid toml [[[`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)

	if err != nil {
		t.Errorf("Load should not error on corrupt file, got: %v", err)
	}
	if cfg == nil {
		t.Error("Load should return config even on corrupt file")
	}

	// Should return defaults
	if cfg.AcceptTimeout != 30*time.Second {
		t.Errorf("Corrupt file should return default timeout, got: %v", cfg.AcceptTimeout)
	}
	if cfg.Identity != "" {
		t.Errorf("Corrupt file should return default identity, got: %s", cfg.Identity)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tmpdir := t.TempDir()
	configDir := filepath.Join(tmpdir, "drift", "config")

	err := EnsureConfigDir(configDir)

	if err != nil {
		t.Errorf("EnsureConfigDir should not error, got: %v", err)
	}

	// Check directory exists
	info, err := os.Stat(configDir)
	if err != nil {
		t.Errorf("Directory should exist after EnsureConfigDir, got: %v", err)
	}
	if !info.IsDir() {
		t.Error("Path should be a directory")
	}

	// Check permissions are 0700
	if info.Mode().Perm() != 0700 {
		t.Errorf("Directory permissions should be 0700, got: %o", info.Mode().Perm())
	}

	// Should not error if called again
	err = EnsureConfigDir(configDir)
	if err != nil {
		t.Errorf("EnsureConfigDir should not error on existing dir, got: %v", err)
	}
}
