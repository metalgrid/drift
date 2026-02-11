package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// Config holds application configuration.
type Config struct {
	DownloadDir   string
	AcceptTimeout time.Duration
	Identity      string
}

// rawConfig is the TOML-decoded structure.
type rawConfig struct {
	DownloadDir   string `toml:"download_dir"`
	AcceptTimeout string `toml:"accept_timeout"`
	Identity      string `toml:"identity"`
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		DownloadDir:   filepath.Join(xdg.UserDirs.Download, "Drift"),
		AcceptTimeout: 30 * time.Second,
		Identity:      "",
	}
}

// DefaultPath returns the default config file path using XDG Base Directory spec.
func DefaultPath() string {
	return filepath.Join(xdg.ConfigHome, "drift", "config.toml")
}

// Load reads a TOML config file and merges with defaults.
// If the file doesn't exist or is corrupt, returns defaults without error.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	// Try to read the file
	data, err := os.ReadFile(path)
	if err != nil {
		// File doesn't exist or can't be read — return defaults silently
		return cfg, nil
	}

	// Try to parse TOML
	var raw rawConfig
	if err := toml.Unmarshal(data, &raw); err != nil {
		// Corrupt TOML — return defaults silently
		return cfg, nil
	}

	// Merge non-empty values from file
	if raw.DownloadDir != "" {
		cfg.DownloadDir = raw.DownloadDir
	}

	if raw.AcceptTimeout != "" {
		duration, err := time.ParseDuration(raw.AcceptTimeout)
		if err == nil {
			cfg.AcceptTimeout = duration
		}
		// If parse fails, keep default
	}

	if raw.Identity != "" {
		cfg.Identity = raw.Identity
	}

	return cfg, nil
}

// EnsureConfigDir creates the config directory with 0700 permissions if it doesn't exist.
func EnsureConfigDir(dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}
