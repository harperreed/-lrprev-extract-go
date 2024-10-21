package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test valid config file
	validConfig := `
input_dir: "/path/to/input"
output_directory: "/path/to/output"
lightroom_db: "/path/to/lightroom.lrcat"
include_size: true
`
	validConfigFile := createTempFile(t, validConfig)
	defer os.Remove(validConfigFile)

	cfg, err := LoadConfig(validConfigFile)
	if err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	if cfg.InputDir != "/path/to/input" {
		t.Errorf("Expected InputDir to be '/path/to/input', got '%s'", cfg.InputDir)
	}
	if cfg.OutputDirectory != "/path/to/output" {
		t.Errorf("Expected OutputDirectory to be '/path/to/output', got '%s'", cfg.OutputDirectory)
	}
	if cfg.LightroomDB != "/path/to/lightroom.lrcat" {
		t.Errorf("Expected LightroomDB to be '/path/to/lightroom.lrcat', got '%s'", cfg.LightroomDB)
	}
	if !cfg.IncludeSize {
		t.Errorf("Expected IncludeSize to be true, got false")
	}

	// Test invalid config file
	invalidConfig := `
input_dir: "/path/to/input"
output_directory: 42  # Should be a string
`
	invalidConfigFile := createTempFile(t, invalidConfig)
	defer os.Remove(invalidConfigFile)

	_, err = LoadConfig(invalidConfigFile)
	if err == nil {
		t.Fatal("Expected error when loading invalid config, got nil")
	}

	// Test non-existent config file
	_, err = LoadConfig("/path/to/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("Expected error when loading non-existent config file, got nil")
	}
}

func createTempFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}
