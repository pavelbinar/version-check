package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// ToolConfig represents the configuration for each tool with a name.
type ToolConfig struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Expect  string `yaml:"expect"`
}

// Config holds an array of tools from the YAML configuration.
type Config struct {
	Tools []ToolConfig `yaml:"tools"`
}

// readConfig reads and parses the config.yaml file.
func readConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// specialCases defines exceptions for version extraction
var specialCases = map[string]string{
	"rsync": `(?:version|protocol)\s+(\d+(?:\.\d+){0,2})`,
}

// extractVersion extracts the version from the command output.
func extractVersion(command string, output string) string {

	if strings.Contains(command, "rsync") {
		// Find the line containing rsync version
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "rsync version") {
				// Use the same regex to parse the version
				re := regexp.MustCompile(`v?(\d+(\.\d+){0,2})`)
				match := re.FindString(line)
				return strings.TrimPrefix(match, "v")
			}
		}
		return ""
	}

	// Default case for other tools
	re := regexp.MustCompile(`v?(\d+(\.\d+){0,2})`)
	match := re.FindString(output)
	return strings.TrimPrefix(match, "v")
}

// compareVersions compares two version strings.
func compareVersions(v1, v2 string) bool {
	// Split versions into parts
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare each part
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] != parts2[i] {
			return false
		}
	}

	// If all parts match and lengths are the same, versions are equal
	return len(parts1) == len(parts2)
}

// checkVersion checks if the output of the command contains the expected version.
func checkVersion(config ToolConfig) bool {
	cmd := exec.Command("sh", "-c", config.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command for %s: %s\n", config.Name, err)
		return false
	}
	actualOutput := strings.TrimSpace(string(output))
	extractedVersion := extractVersion(config.Command, actualOutput)
	if extractedVersion == "" {
		fmt.Printf("%s: No version found in output\n", config.Name)
		fmt.Printf("Command output: %s\n", actualOutput)
		return false
	}
	if compareVersions(extractedVersion, config.Expect) {
		return true
	} else {
		fmt.Printf("%s version mismatch: Expected '%s', got '%s'\n", config.Name, config.Expect, extractedVersion)
		fmt.Printf("Command output: %s\n", actualOutput)
		return false
	}
}

func main() {
	config, err := readConfig()
	if err != nil {
		fmt.Printf("Failed to read config.yaml: %s\n", err)
		os.Exit(1)
	}

	allMatch := true
	for _, toolConfig := range config.Tools {
		if !checkVersion(toolConfig) {
			allMatch = false
		}
	}

	if allMatch {
		fmt.Println("Versions OK")
	} else {
		os.Exit(1) // Exit with error status if any version does not match
	}
}
