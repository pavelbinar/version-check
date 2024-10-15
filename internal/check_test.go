package check

import (
	"os"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	// Create a temporary test config file with various version formats
	tempConfig := []byte(`
tools:
  - name: "FullVersion"
    command: "echo 'FullVersion 1.22.3'"
    expect: "1.22.3"
  - name: "ShortVersion"
    command: "echo 'ShortVersion v1.22'"
    expect: "1.22"
  - name: "SingleDigit"
    command: "echo 'SingleDigit version 2'"
    expect: "2"
  - name: "VersionMismatch"
    command: "echo 'VersionMismatch 1.22.4'"
    expect: "1.22.3"
  - name: "PrefixVersion"
    command: "echo 'PrefixVersion v2.0.1'"
    expect: "2.0.1"
  - name: "SuffixVersion"
    command: "echo 'SuffixVersion 3.1.4-alpha'"
    expect: "3.1.4"
  - name: "ErrorCommand"
    command: "non_existent_command"
    expect: "1.0.0"
`)

	tmpfile, err := os.CreateTemp("", "test_config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(tempConfig); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Temporarily replace the config file path
	originalConfig := "config.yaml"
	os.Rename(originalConfig, originalConfig+".bak")
	os.Rename(tmpfile.Name(), originalConfig)
	defer func() {
		os.Rename(originalConfig, tmpfile.Name())
		os.Rename(originalConfig+".bak", originalConfig)
	}()

	// Run the tests
	config, err := readTestConfig(originalConfig)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	tests := []struct {
		name     string
		tool     ToolConfig
		expected bool
	}{
		{"Full version match", config.Tools[0], true},
		{"Short version match", config.Tools[1], true},
		{"Single digit version match", config.Tools[2], true},
		{"Version mismatch", config.Tools[3], false},
		{"Prefix version match", config.Tools[4], true},
		{"Suffix version match", config.Tools[5], true},
		{"Command error", config.Tools[6], false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkVersion(tt.tool)
			if result != tt.expected {
				t.Errorf("checkVersion(%v) = %v, want %v", tt.tool, result, tt.expected)
			}
		})
	}
}

func TestVersionExtraction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Full version", "Tool version 1.22.3", "1.22.3"},
		{"Short version", "Tool v1.22", "1.22"},
		{"Single digit", "Tool version 2", "2"},
		{"Version with prefix", "Tool v2.0.1", "2.0.1"},
		{"Version with suffix", "Tool 3.1.4-alpha", "3.1.4"},
		{"Complex output", "Tool (v4.5.6) [built 2023-05-01]", "4.5.6"},
		{"No version", "Tool has no version", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := ToolConfig{
				Name:    "TestTool",
				Command: "echo '" + tt.input + "'",
				Expect:  tt.expected,
			}
			result := checkVersion(tool)
			if tt.expected == "" {
				if result {
					t.Errorf("Expected false for input without version, got true")
				}
			} else {
				if !result {
					t.Errorf("Expected true for input '%s', got false", tt.input)
				}
			}
		})
	}
}

func TestExtractVersionWithSpecialCase(t *testing.T) {
	tests := []struct {
		name              string
		toolName          string
		commandOutput     string
		expectedVersion   string
		expectedExtracted bool
	}{
		{
			name:              "Rsync special case",
			toolName:          "rsync",
			commandOutput:     "openrsync: protocol version 29\nrsync version 2.6.9 compatible",
			expectedVersion:   "2.6.9",
			expectedExtracted: true,
		},
		{
			name:              "Normal case",
			toolName:          "normalTool",
			commandOutput:     "normalTool version v1.2.3",
			expectedVersion:   "1.2.3",
			expectedExtracted: true,
		},
		{
			name:              "No version in output",
			toolName:          "noVersionTool",
			commandOutput:     "This output contains no version information",
			expectedVersion:   "",
			expectedExtracted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractedVersion := extractVersion(tt.toolName, tt.commandOutput)
			if extractedVersion != tt.expectedVersion {
				t.Errorf("Expected version %s, but got %s", tt.expectedVersion, extractedVersion)
			}

			// Test checkVersion function
			tool := ToolConfig{
				Name:    tt.toolName,
				Command: "echo '" + tt.commandOutput + "'",
				Expect:  tt.expectedVersion,
			}
			result := checkVersion(tool)
			if result != tt.expectedExtracted {
				t.Errorf("checkVersion() = %v, want %v", result, tt.expectedExtracted)
			}
		})
	}
}

// Rename this function to avoid conflict
func readTestConfig(filename string) (*Config, error) {
	return readConfig(filename)
}
